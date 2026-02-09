<template>
  <div class="flex h-screen">
    <!-- 左侧会话列表 -->
    <div class="w-60 border-r border-gray-200 bg-white flex flex-col">
      <div class="p-4 border-b border-gray-200">
        <h2 class="font-bold text-gray-800">会话列表</h2>
        <p class="text-sm text-gray-500">{{ activeSessions.length }} 个进行中</p>
      </div>
      
      <div class="flex-1 overflow-y-auto">
        <div
          v-for="session in sessions"
          :key="session.id"
          @click="selectSession(session)"
          class="p-4 border-b border-gray-100 cursor-pointer transition-colors"
          :class="selectedSession?.id === session.id ? 'bg-blue-50 border-l-4 border-l-blue-600' : 'hover:bg-gray-50'"
        >
          <div class="flex items-start gap-3">
            <el-avatar size="large" :src="session.avatar" class="w-10 h-10">
                {{ session.name?.charAt(0) || 'U' }}
            </el-avatar>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <span class="font-medium text-gray-800">{{ session.name }}</span>
                <span class="text-xs text-gray-500">{{ session.time }}</span>
              </div>
              <p class="text-sm text-gray-600 truncate">{{ session.lastMessage }}</p>
              <div class="flex items-center gap-2 mt-1">
                <el-tag :type="session.status === '进行中' ? 'success' : 'info'" size="small">
                  {{ session.status }}
                </el-tag>
                <span v-if="session.unread > 0" class="text-xs text-red-500 font-medium">
                  {{ session.unread }} 条未读
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 中间聊天区域 -->
    <div class="flex-1 flex flex-col bg-gray-50">
      <div v-if="selectedSession" class="flex-1 flex flex-col">
        <!-- 聊天头部 -->
        <div class="p-4 bg-white border-b border-gray-200 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <el-avatar size="large" :src="selectedSession.avatar" class="w-10 h-10">
                {{ selectedSession.name?.charAt(0) || 'U' }}
            </el-avatar>
            <div>
              <p class="font-medium text-gray-800">{{ selectedSession.name }}</p>
              <p class="text-xs text-gray-500">{{ selectedSession.email }}</p>
            </div>
          </div>
          <div class="flex gap-2">
            <el-button type="primary" size="small">转接</el-button>
            <el-button type="danger" size="small">结束会话</el-button>
          </div>
        </div>

        <!-- 消息列表 -->
        <div class="flex-1 overflow-y-auto p-4 space-y-4">
          <div
            v-for="message in messages"
            :key="message.id"
            class="flex"
            :class="message.isSelf ? 'justify-end' : 'justify-start'"
          >
            <div
              class="max-w-xs px-4 py-2 rounded-lg"
              :class="message.isSelf ? 'bg-blue-600 text-white' : 'bg-white text-gray-800'"
            >
              <p>{{ message.content }}</p>
              <p class="text-xs mt-1 opacity-70">{{ message.time }}</p>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="p-4 bg-white border-t border-gray-200">
            <div class="flex items-center gap-3 mb-3">
                <button v-for="item in toolbarItems" :key="item.action" class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 cursor-pointer" @click="handleToolbarClick(item.action)" :title="item.label">
                    <div class="w-6 h-6 text-gray-600" v-html="item.svg"></div>
                </button>
            </div>
            <div class="relative w-full">
                <textarea v-model="inputMessage" rows="8" placeholder="输入消息..." @keydown.enter.exact.prevent="onEnter" @keydown.ctrl.enter.exact.prevent="sendMessage" class="w-full px-3 py-2 pr-28  focus:ring-1 focus:ring-blue-500 outline-none  text-sm"></textarea>
                <div class=" absolute right-2 bottom-4 flex items-center">
                    <button type="button" @click="sendMessage" class="bg-blue-600 hover:bg-blue-700 text-white text-base pl-2 py-1 pr-2 rounded-l-md rounded-r-none flex items-center border-r border-blue-300"><el-icon class="text-white text-xl! mr-1"><Promotion /></el-icon>发送</button>
                    <button type="button" @click="toggleMenu" class="bg-blue-600 hover:bg-blue-700 text-white text-base py-1 px-2 rounded-r-md rounded-l-none"><el-icon><ArrowDown /></el-icon></button>
                    <div v-if="showMenu" ref="menuRef" class="absolute bottom-full right-0 mb-1 bg-white border border-gray-200 rounded-md shadow-lg z-10 w-54">
                    <div @click="setSendMode('enter')" class="px-3 py-2 text-sm cursor-pointer hover:bg-gray-100 flex items-center" :class="{ 'text-blue-600 font-medium': sendMode === 'enter' }"><span :class="['mr-2',sendMode !== 'enter'?'ml-2':'']"> {{ sendMode === 'enter'?'✓':'' }}</span>按 Enter 键发送消息</div>
                    <div @click="setSendMode('ctrl-enter')" class="px-3 py-2 text-sm cursor-pointer hover:bg-gray-100 flex items-center" :class="{ 'text-blue-600 font-medium': sendMode === 'ctrl-enter' }"><span :class="['mr-2',sendMode !== 'ctrl-enter'?'ml-2':'']">{{ sendMode === 'ctrl-enter'?'✓':'' }}</span>按 Ctrl + Enter 键发送消息</div>
                    </div>
                </div>
            </div>
        </div>


      </div>

      <!-- 空状态 -->
      <div v-else class="flex-1 flex items-center justify-center text-gray-400">
        <div class="text-center">
          <el-icon :size="64" class="mb-4"><ChatDotRound /></el-icon>
          <p>选择一个会话开始聊天</p>
        </div>
      </div>
    </div>

    <!-- 右侧访客资料 -->
    <div v-if="selectedSession" class="w-80 border-l border-gray-200 bg-white flex flex-col">
      <div class="px-4  py-6 border-b border-gray-200">
        <h2 class="font-bold text-gray-800">访客资料</h2>
      </div>
      
      <div class="flex-1 overflow-y-auto p-4">
        <div class="flex flex-col items-center gap-4">
          <el-avatar :size="80" :src="selectedSession.avatar">
            {{ selectedSession.name?.charAt(0) }}
          </el-avatar>
          <h3 class="text-xl font-bold text-gray-800">{{ selectedSession.name }}</h3>
        </div>
        
        <div class="space-y-4">
          <div>
            <p class="text-sm text-gray-500 mb-1">邮箱</p>
            <p class="text-gray-800 font-medium">{{ selectedSession.email }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 mb-1">IP地址</p>
            <p class="text-gray-800 font-medium">{{ selectedSession.ip || '192.168.1.1' }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 mb-1">浏览器</p>
            <p class="text-gray-800 font-medium">{{ selectedSession.browser || 'Chrome 120' }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 mb-1">操作系统</p>
            <p class="text-gray-800 font-medium">{{ selectedSession.os || 'Windows 11' }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 mb-1">访问次数</p>
            <p class="text-gray-800 font-medium">{{ selectedSession.visits || 5 }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 mb-1">最后访问</p>
            <p class="text-gray-800 font-medium">{{ selectedSession.lastVisit || '10:30' }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ChatDotRound, Promotion, ArrowDown } from '@element-plus/icons-vue'
import { Emoji, Picture, Folder, Audio, Block, Transfer, History, Robot, Snippets } from '../script/svg.js'

const selectedSession = ref(null)
const inputMessage = ref('')
const sendMode = ref('ctrl-enter')
const showMenu = ref(false)
const menuRef = ref(null)

// ✅ 定义工具栏配置（使用 svg.js 中的完整 SVG 图标）
const toolbarItems = [
  { action: 'emoji',     label: '表情',     svg: Emoji() },
  { action: 'image',     label: '图片',     svg: Picture() },
  { action: 'folder',    label: '文件',     svg: Folder() },
  { action: 'voice',     label: '语音',     svg: Audio() },
  { action: 'block',     label: '拉黑',     svg: Block() },
  { action: 'transfer',  label: '转接',     svg: Transfer() },
  { action: 'history',   label: '历史记录', svg: History() },
  { action: 'robot',     label: '机器人',   svg: Robot() },
  { action: 'snippets',  label: '快捷回复', svg: Snippets() }
];

const sessions = ref([
  { id: 1, name: '访客A', email: 'visitor1@example.com', avatar: '', status: '进行中', time: '10:30', lastMessage: '你好，请问有什么可以帮您的吗？', unread: 0, ip: '192.168.1.100', browser: 'Chrome 120', os: 'Windows 11', visits: 12, lastVisit: '10:30' },
  { id: 2, name: '访客B', email: 'visitor2@example.com', avatar: '', status: '进行中', time: '10:25', lastMessage: '我想咨询一下产品价格', unread: 2, ip: '192.168.1.101', browser: 'Firefox 115', os: 'macOS 14', visits: 8, lastVisit: '09:45' },
  { id: 3, name: '访客C', email: 'visitor3@example.com', avatar: '', status: '进行中', time: '10:20', lastMessage: '好的，谢谢', unread: 0, ip: '192.168.1.102', browser: 'Safari 15', os: 'iOS 16', visits: 15, lastVisit: '10:20' },
  { id: 4, name: '访客D', email: 'visitor4@example.com', avatar: '', status: '已结束', time: '10:15', lastMessage: '再见', unread: 0, ip: '192.168.1.103', browser: 'Edge 120', os: 'Windows 10', visits: 3, lastVisit: '08:30' },
  { id: 5, name: '访客E', email: 'visitor5@example.com', avatar: '', status: '进行中', time: '10:10', lastMessage: '请问有优惠活动吗？', unread: 1, ip: '192.168.1.104', browser: 'Chrome 119', os: 'Android 13', visits: 6, lastVisit: '10:10' }
])

const activeSessions = computed(() => sessions.value.filter(s => s.status === '进行中'))

const messages = ref([])

const selectSession = (session) => {
  selectedSession.value = session
  messages.value = [
    { id: 1, content: '你好，请问有什么可以帮您的吗？', time: '10:30', isSelf: false },
    { id: 2, content: '我想咨询一下产品价格', time: '10:31', isSelf: true },
    { id: 3, content: '好的，请问您想了解哪个产品？', time: '10:32', isSelf: false }
  ]
}

const sendMessage = () => {
  if (!inputMessage.value.trim()) return
  
  messages.value.push({
    id: messages.value.length + 1,
    content: inputMessage.value,
    time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
    isSelf: true
  })
  
  inputMessage.value = ''
}

const toggleMenu = () => (showMenu.value = !showMenu.value)
const setSendMode = (mode) => {
  sendMode.value = mode
  showMenu.value = false
}
</script>
