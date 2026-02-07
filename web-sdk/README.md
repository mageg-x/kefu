当然可以！下面我将你提出的几个关键问题——**SDK 接入、用户身份升级、客服转接**——整合成一个**完整的轻量级网页客服系统架构流程**，完全基于 **WuKongIM 作为消息通道**，并最大化复用其能力（包括消息存储），同时只引入**最小必要的业务元数据**。

---

## 🎯 系统目标

- 访客通过 `<script>` 嵌入聊天窗口（支持匿名 & 已登录）
- 客服在 `web-admin` 后台处理会话
- 支持：  
  ✅ 匿名访客 → 登录后身份升级（消息连续）  
  ✅ 智能客服 ↔ 人工客服 转接（用户无感）  
  ✅ 所有历史消息由 WuKongIM 存储，不重复存消息体  
  ✅ 企业可透传自有用户 ID 给客服

---

## 🧱 整体架构

```
+------------------+       +------------------+       +------------------+
|   用户网站        |       |     你的 Server   |       |   WuKongIM Server |
| (嵌入 widget.min.js)| <---> | (API + 逻辑控制)  | <---> | (消息/用户/频道)  |
+------------------+       +------------------+       +------------------+
         ↑
         | WebSocket / HTTP
+------------------+
|   web-admin      |
| (客服后台)        |
+------------------+
```

---

## 🔑 核心设计原则

1. **会话（Session）是核心单位**，不是用户
2. **使用自定义频道（Custom Channel）**，而非 P2P，以支持转接
3. **消息存储完全依赖 WuKongIM**
4. **仅存储轻量元数据**：会话 ↔ 频道 ↔ 用户身份 映射
5. **访客身份可动态升级**，不影响已有消息

---

## 📦 一、前端 SDK 接入（web-sdk）

### 1. 嵌入方式（不变）
```html
<script src="./dist/widget.min.js" data-kefu-appid="your_app_id" defer></script>
```

### 2. 初始化流程（web-sdk 内部）

```js
// 1. 生成或复用 sessionId（持久化）
const sessionId = localStorage.getItem('kefu_session_id') || 'sess_' + uuid();

// 2. 调用你的 server 初始化会话
const res = await fetch('/api/init-session', {
  method: 'POST',
  body: JSON.stringify({ appId: 'your_app_id', sessionId })
});

const { channelID, wukongToken, visitorKey } = res.json();

// 3. 初始化 WuKongIM Web SDK
wukongIM.connect(wukongToken);
wukongIM.joinChannel(channelID); // 自定义频道

// 4. 暴露 API 给企业
window.KefuWidget = {
  setUserInfo(userInfo) {
    // 上报真实身份（如用户登录后）
    fetch('/api/bind-user', { 
      method: 'POST',
      body: JSON.stringify({ sessionId, ...userInfo }) 
    });
  }
};
```

> ✅ 此时访客已进入频道，可收发消息。

---

## 🔄 二、用户身份升级（匿名 → 实名）

### 场景
用户先匿名咨询，后登录，希望客服看到真实信息。

### 流程

1. **企业调用**（登录成功后）：
   ```js
   KefuWidget.setUserInfo({ userId: "user_123", name: "张三", avatar: "..." });
   ```

2. **Server 处理 `/api/bind-user`**：
   - 查 `kefu_sessions` 表，找到 `sessionId` 对应的记录
   - 更新 `real_user_id`, `display_name` 等字段
   - **调用 WuKongIM 发送系统消息到频道**：
     ```json
     {
       "type": "visitor_identified",
       "data": { "name": "张三", "userId": "user_123" }
     }
     ```

3. **客户端表现**：
   - **访客端**：无变化（继续聊天）
   - **客服端**：收到系统消息，更新访客头像/昵称，后续消息显示“张三”

> ✅ 历史消息无需迁移 —— 因为所有消息都在**同一个频道**中。

---

## 🔄 三、客服转接（智能 → 人工 / A → B）

### 前提
- 会话使用 **自定义频道**（如 `channelID = "kefu_sess_a1b2"`）
- 初始成员：`[guest_xyz, agent_robot]`

### 转接流程

1. **发起转接**（客服点击“转接”按钮）：
   ```http
   POST /api/transfer-session
   {
     "sessionId": "sess_a1b2",
     "toAgentId": "agent_zhangsan",
     "reason": "用户请求人工"
   }
   ```

2. **Server 执行**：
   - 验证权限
   - 调用 WuKongIM API：
     ```go
     wukong.AddChannelMember("kefu_sess_a1b2", "agent_zhangsan");
     // 可选：wukong.RemoveChannelMember(..., "agent_robot");
     ```
   - 发送系统消息到频道：
     ```json
     {
       "type": "transfer",
       "from": "智能客服",
       "to": "张三（人工客服）"
     }
     ```
   - 更新 `kefu_sessions.current_agent = "agent_zhangsan"`

3. **客户端表现**：
   - **访客**：看到“已转接至人工客服 张三”，继续在同一窗口聊天
   - **新客服**：收到通知，打开会话后加载该频道**全部历史消息**
   - **原客服**：界面标记“已转出”

> ✅ 消息连续、无副本、无断裂。

---

## 🗃 四、最小业务数据库设计（仅元数据）

### 表 1：会话主表
```sql
CREATE TABLE kefu_sessions (
  id BIGSERIAL PRIMARY KEY,
  session_key TEXT NOT NULL UNIQUE,     -- 来自前端 localStorage
  visitor_key TEXT NOT NULL,            -- visitor_xxx（逻辑身份）
  channel_id TEXT NOT NULL,             -- WuKongIM 频道 ID
  real_user_id TEXT,                    -- 企业用户 ID（可为空）
  display_name TEXT,
  avatar_url TEXT,
  current_agent TEXT NOT NULL,          -- 当前负责客服
  status TEXT DEFAULT 'active',         -- active / closed
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);
```

### 表 2：转接日志（可选）
```sql
CREATE TABLE kefu_transfers (
  id BIGSERIAL PRIMARY KEY,
  session_key TEXT,
  from_agent TEXT,
  to_agent TEXT,
  reason TEXT,
  transferred_at TIMESTAMPTZ
);
```

> 💡 **不存储任何消息内容！消息全由 WuKongIM 管理。**

---

## 🧩 五、WuKongIM 关键使用点

| 功能 | 使用方式 |
|------|--------|
| 用户认证 | Server 动态创建临时访客用户 + 签发 token |
| 消息通道 | **自定义频道（channelType=100）** |
| 历史消息 | `getHistoryMessages(channelID)` |
| 成员管理 | `Add/RemoveChannelMember`（用于转接） |
| 用户信息 | 通过系统消息传递身份变更，不依赖 WuKongIM user metadata |

---

## ✅ 六、优势总结

| 能力 | 是否支持 | 说明 |
|------|--------|------|
| 匿名访客接入 | ✅ | 临时 userID + token |
| 登录后身份升级 | ✅ | 通过系统消息通知，不换频道 |
| 消息连续性 | ✅ | 所有消息在同一自定义频道 |
| 客服转接 | ✅ | 动态变更频道成员 |
| 不重复存消息 | ✅ | 完全依赖 WuKongIM 存储 |
| 企业透传用户 ID | ✅ | 通过 `setUserInfo` 上报 |
| 客服看到完整上下文 | ✅ | 加载频道全部历史 |

---

## 🚀 下一步建议

1. **在 server 中实现 `/api/init-session`**  
   → 分配客服、创建频道、返回 token + channelID

2. **web-sdk 改为加入自定义频道**（而非 P2P）

3. **web-admin 改为按 `channelID` 加载会话**

4. **提供 `KefuWidget.setUserInfo` 接口文档给企业**

