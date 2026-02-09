package models

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SessionStatusUnAssigned = "unassigned" // 未分配客服
	SessionStatusUnRead     = "unread"     // 未读新消息
	SessionStatusUnReply    = "unreply"    // 有未回复
	SessionStatusAssigned   = "assigned"   // 已分配客服
	SessionStatusFollowUP   = "follow"     // 待跟进
	SessionStatusClosed     = "closed"     // 会话结束
)

type Session struct {
	SID                string `json:"sid"`                    // m:{visitor_id}:{app_id}:{session_seq}
	CurAgentID         string `json:"cur_agent_id,omitempty"` // 当前负责客服
	CreatedAt          int64  `json:"created_at"`             // 创建时间
	LastVisitorMsgTime int64  `json:"last_visitor_msg_time"`  // 最后访客消息时间
	LastAgentReplyTime int64  `json:"last_agent_reply_time"`  // 最后客服回复消息时间
	LastAgentReadTime  int64  `json:"last_agent_read_time"`   // 最后客服已读消息时间
	Closed             bool   `json:"closed"`                 // 会话是否关闭
	FollowUp           bool   `json:"need_follow_up"`         // 会话是否需要跟进
}

// 由 session_id 提取 visitor_id, app_id, session_seq
func ParseSessionID(sessionID string) (visitorID, appID string, sessionSeq uint32) {
	s := Session{SID: sessionID}
	return s.ParseSid()
}

func GetSessionID(visitorID, appID string, sessionSeq uint32) string {
	return fmt.Sprintf("m:%s:%s:%010d", visitorID, appID, sessionSeq)
}

// 1. 访客发消息
func (s *Session) OnVisitorMessage(ts int64) {
	if s.Closed {
		return
	}
	s.LastVisitorMsgTime = ts
}

// 2. 客服发送普通回复
func (s *Session) OnAgentReply(ts int64) {
	if s.Closed {
		return
	}
	s.LastAgentReplyTime = ts
	s.LastAgentReadTime = ts // 回复即视为已读
}

// 3. 客服接单
func (s *Session) AssignAgent(agentID string, ts int64) {
	if s.Closed {
		return
	}
	s.CurAgentID = agentID
	// 注意：不改 Read/Reply 时间！
}

// 4. 客服点开会话（触发 badge 清除）
func (s *Session) MarkRead(ts int64) {
	if s.Closed {
		return
	}
	if s.CurAgentID != "" {
		s.LastAgentReadTime = ts
	}
}

// 5. 客服标记“需后续跟进”
func (s *Session) MarkFollowUp() {
	if s.Closed {
		return
	}
	s.FollowUp = true
	// 关键：清空已读时间，确保 badge 出现
	s.LastAgentReadTime = 0
}

// 6. 关闭会话
func (s *Session) Close() {
	s.Closed = true
	s.FollowUp = false
}

// 7. 获取会话状态
func (s *Session) Status() string {
	if s.Closed {
		return SessionStatusClosed
	}
	if s.CurAgentID == "" {
		return SessionStatusUnAssigned
	}
	if s.LastVisitorMsgTime > s.LastAgentReadTime {
		return SessionStatusUnRead
	}
	if s.LastVisitorMsgTime > s.LastAgentReplyTime {
		return SessionStatusUnReply
	}
	if s.FollowUp {
		return SessionStatusFollowUP
	}
	return SessionStatusAssigned
}

// 8、由 sid 提取 visitor_id, app_id, session_seq
func (s *Session) ParseSid() (visitorID, appID string, sessionSeq uint32) {
	parts := strings.Split(s.SID, ":")
	if len(parts) != 4 {
		return "", "", 0
	}

	visitorID = parts[1]
	appID = parts[2]
	sessionSeq64, err := strconv.ParseUint(parts[3], 10, 64)
	if err != nil {
		return "", "", 0
	}
	sessionSeq = uint32(sessionSeq64)
	return
}

// 9. 获取 visitor_id
func (s *Session) VisitorID() string {
	visitor, _, _ := s.ParseSid()
	return visitor
}

// 10. 获取 app_id
func (s *Session) AppID() string {
	_, appID, _ := s.ParseSid()
	return appID
}

// 11. 获取 agent_id
func (s *Session) AgentID() string {
	return s.CurAgentID
}

// 12. 获取 session_seq
func (s *Session) SessionSeq() uint32 {
	_, _, sessionSeq := s.ParseSid()
	return sessionSeq
}
