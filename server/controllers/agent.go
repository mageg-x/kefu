package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"

	"kefu-server/models"
	"kefu-server/service"
	"kefu-server/utils/logger"
)

const (
	MessageTypeRsp = "message.rsp"
)

// AgentConn 表示一个客服的 WebSocket 连接
type AgentConn struct {
	Conn     *websocket.Conn
	AgentID  string
	SendChan chan []byte
	Done     chan struct{}
}

// 全局客服连接池：agent_id => *AgentConn
var (
	agentConns sync.Map // map[string]*AgentConn
)

// 注册客服连接
func registerAgentConn(agentID string, conn *AgentConn) {
	agentConns.Store(agentID, conn)
}

// 注销客服连接
func unregisterAgentConn(agentID string) {
	agentConns.Delete(agentID)
}

// 向客服推送消息（供系统调用）
func PushMessageToAgent(agentID, sessionID string, msg *models.Message) {
	if v, ok := agentConns.Load(agentID); ok {
		conn := v.(*AgentConn)
		payload, _ := json.Marshal(struct {
			Type      string          `json:"type"`
			SessionID string          `json:"session_id"`
			Message   *models.Message `json:"message"`
		}{
			Type:      "message.req",
			SessionID: sessionID,
			Message:   msg,
		})

		select {
		case conn.SendChan <- payload:
		default:
			logger.Warnf("Agent %s send buffer full", agentID)
		}
	}
}

// AgentController 客服控制器
type AgentController struct{}

// WSHandler 处理客服 WebSocket 连接
func (ac *AgentController) WSHandler(c *gin.Context) {
	// 从上下文中获取认证后的用户信息（由 auth 中间件注入）
	userName, exists := c.Get("userName") // 假设 userName 是 agent_id
	if !exists || userName == "" {
		logger.Errorf("Agent ID not found in context")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	agentID := userName.(string)

	// 可选：验证 agent 是否存在
	us := service.GetUserService()
	if us == nil {
		logger.Errorf("User service not initialized")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 进一步校验 agent 的合法性
	agent, err := us.GetUser(agentID)
	if agent == nil || err != nil {
		logger.Errorf("Agent %s does not exist", agentID)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if agent.Role != "agent" || agent.Active != true {
		logger.Errorf("Agent %s does not have agent role or is not active", agentID)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 升级 WebSocket
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		logger.Errorf("WebSocket upgrade failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer conn.CloseNow()

	// 创建连接对象
	agentConn := &AgentConn{
		Conn:     conn,
		AgentID:  agentID,
		SendChan: make(chan []byte, 256),
		Done:     make(chan struct{}),
	}

	// 注册到连接池
	registerAgentConn(agentID, agentConn)
	defer unregisterAgentConn(agentID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ac.readLoop(ctx, agentConn)
	go ac.writeLoop(ctx, agentConn)

	logger.Infof("Agent connected: %s", agentID)
	<-agentConn.Done
	logger.Infof("Agent disconnected: %s", agentID)
}

func (ac *AgentController) readLoop(ctx context.Context, conn *AgentConn) {
	defer close(conn.Done)

	for {
		readCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		_, data, err := conn.Conn.Read(readCtx)
		cancel()

		if err != nil {
			if websocket.CloseStatus(err) == -1 {
				logger.Errorf("Agent %s read error: %v", conn.AgentID, err)
			}
			return
		}

		var req struct {
			Type    string `json:"type"`
			Session string `json:"session_id"`
			Payload string `json:"payload"`
		}
		if err := json.Unmarshal(data, &req); err != nil {
			continue
		}

		if req.Session == "" {
			continue
		}

		ac.handleMessage(conn.AgentID, req.Session, req.Type, req.Payload)
	}
}

func (ac *AgentController) writeLoop(ctx context.Context, conn *AgentConn) {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-conn.Done:
			return
		case msg := <-conn.SendChan:
			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := conn.Conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				logger.Errorf("Write to agent %s failed: %v", conn.AgentID, err)
				return
			}
		case <-ticker.C:
			conn.Conn.Ping(ctx)
		}
	}
}

func (ac *AgentController) handleMessage(agentID, sessionID, actionType, payload string) {
	ss := service.GetSessionService()
	session, err := ss.GetSession(sessionID)
	if err != nil || session == nil {
		logger.Errorf("Session %s does not exist", sessionID)
		return
	}

	// 权限校验
	if session.CurAgentID != agentID {
		logger.Errorf("Agent %s not assigned to session %s", agentID, sessionID)
		return
	}

	now := time.Now().Unix()

	switch actionType {
	case "message.rsp":
		// 保存客服回复
		ms := service.GetMsgService()
		if ms == nil { // 单例
			logger.Errorf("msg service is not initialized")
			return
		}
		msg := models.Message{
			Content:   payload,
			MsgType:   "message.rsp",
			Timestamp: now,
		}
		msgID, err := ms.SaveMessage(session.VisitorID(), session.AppID(), session.SessionSeq(), &msg)
		if err != nil {
			logger.Errorf("Save message failed: %v", err)
		}
		msg.MsgID = msgID

		// 更新会话状态
		session.OnAgentReply(now)
		ss.SaveSession(session)

		PushMessageToVisitor(session.VisitorID(), sessionID, &msg)

	case "close_session":
		session.Close()
		ss.SaveSession(session)

	case "mark_follow_up":
		session.MarkFollowUp()
		ss.SaveSession(session)

	default:
		logger.Debugf("Unhandled agent action: %s", actionType)
	}
}
