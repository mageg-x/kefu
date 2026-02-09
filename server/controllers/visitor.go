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
	MessageTypeReq    = "message.req"
	MessageTypeTyping = "message.typing"
	MessageTypeClose  = "message.close"
)

// VisitorConn 封装访客连接
type VisitorConn struct {
	Conn      *websocket.Conn
	SessionID string
	SendChan  chan []byte
	Done      chan struct{}
}

// 全局访客连接池：session_id => *VisitorConn
var (
	visitorConns = make(map[string]*VisitorConn)
	visitorMu    sync.RWMutex
)

// 注册连接
func registerVisitorConn(sessionID string, conn *VisitorConn) {
	visitorMu.Lock()
	visitorConns[sessionID] = conn
	visitorMu.Unlock()
}

// 注销连接
func unregisterVisitorConn(sessionID string) {
	visitorMu.Lock()
	delete(visitorConns, sessionID)
	visitorMu.Unlock()
}

// 推送消息给访客（供客服系统调用）
func PushMessageToVisitor(visitorID, sessionID string, msg *models.Message) error {
	visitorMu.RLock()
	conn, ok := visitorConns[sessionID]
	visitorMu.RUnlock()

	if !ok || conn == nil {
		logger.Errorf("Visitor %s not found", sessionID)
		return nil // 访客不在线，静默丢弃（或可存离线消息）
	}

	payload, _ := json.Marshal(struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}{
		Type:    msg.MsgType,
		Payload: msg.Content,
	})

	select {
	case conn.SendChan <- payload:
	default:
		logger.Warnf("Visitor %s send buffer full", sessionID)
	}
	return nil
}

type VisitorController struct{}

func (vc *VisitorController) WSHandler(c *gin.Context) {
	visitorID := c.Query("visitor_id")
	appID := c.Query("app_id")

	if visitorID == "" || appID == "" {
		logger.Errorf("Visitor ID or App ID not found")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !vc.isValidOrigin(appID, c.GetHeader("Origin"), c.GetHeader("Referer")) {
		logger.Errorf("Origin not allowed %s", appID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ss := service.GetSessionService()
	if ss == nil {
		logger.Errorf("Session service not initialized")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	session, err := ss.GetOrCreateSession(visitorID, appID)
	if err != nil {
		logger.Errorf("Failed to get session: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		logger.Errorf("Failed to accept websocket connection: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer conn.CloseNow()

	// 创建连接对象
	visitorConn := &VisitorConn{
		Conn:      conn,
		SessionID: session.SID,
		SendChan:  make(chan []byte, 128),
		Done:      make(chan struct{}),
	}

	// 注册到连接池
	registerVisitorConn(session.SID, visitorConn)
	defer unregisterVisitorConn(session.SID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go vc.readLoop(ctx, visitorConn)
	go vc.writeLoop(ctx, visitorConn)

	<-visitorConn.Done
	logger.Infof("Visitor disconnected: %s", session.SID)
}

func (vc *VisitorController) readLoop(ctx context.Context, vconn *VisitorConn) {
	defer close(vconn.Done)

	for {
		readCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		_, data, err := vconn.Conn.Read(readCtx)
		cancel()

		if err != nil {
			if websocket.CloseStatus(err) == -1 {
				logger.Errorf("Read error: %v", err)
			}
			return
		}

		var req struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(data, &req); err != nil {
			continue
		}

		vc.handleMessage(vconn.SessionID, req.Type, string(req.Payload))
	}
}

func (vc *VisitorController) writeLoop(ctx context.Context, vconn *VisitorConn) {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-vconn.Done:
			return
		case msg := <-vconn.SendChan:
			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := vconn.Conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				logger.Errorf("Write error: %v", err)
				return
			}
		case <-ticker.C:
			vconn.Conn.Ping(ctx)
		}
	}
}

func (vc *VisitorController) handleMessage(sessionID, msgType, content string) {
	ss := service.GetSessionService()
	if ss == nil { // 单例
		logger.Errorf("Session service not initialized")
		return
	}
	session, err := ss.GetSession(sessionID)
	if err != nil || session == nil {
		logger.Errorf("Session %s does not exist", sessionID)
		return
	}

	now := time.Now().Unix()
	session.OnVisitorMessage(now)

	ms := service.GetMsgService()
	if ms == nil {
		logger.Errorf("Message service not initialized")
		return
	}

	msg := models.Message{Content: content, MsgType: msgType, Timestamp: now}
	msgID, err := ms.SaveMessage(session.VisitorID(), session.AppID(), session.SessionSeq(), &msg)
	if err != nil {
		logger.Errorf("Failed to save message: %v", err)
		return
	}
	msg.MsgID = msgID

	ss.SaveSession(session)

	// 自动分配客服
	if session.CurAgentID == "" {
		us := service.GetUserService()
		if agent, _ := us.FindAgent(session.AppID()); agent != nil {
			session.AssignAgent(agent.Username, now)
			ss.SaveSession(session)
		}
	}

	// 在分配客服后，推送消息给客服
	if session.CurAgentID != "" {
		PushMessageToAgent(session.AgentID(), session.SID, &msg)
	}
}

func (vc *VisitorController) isValidOrigin(appID, origin, referer string) bool {
	app := models.GetApp(appID)
	if app == nil {
		return false
	}
	return models.IsDomainAllowed(origin, referer, app.AllowDomain)
}
