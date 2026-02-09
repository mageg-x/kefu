package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"

	"kefu-server/models"
	"kefu-server/store"
	"kefu-server/utils/logger"
)

// 1. 会话存储
// s:{visitor_id}:{app_id}:{session_seq} 为 会话 id
// s:alice:shop123:1234
// │   │     │       │
// │   │     │       └─ session seq
// │   │     └─ app_id
// │   └─ visitor_id

type SessionService struct {
	kv *badger.DB
}

const (
	SessionTimeout = 24 * time.Hour // 可配置
)

var (
	instSessionService *SessionService
)

func GetSessionService() *SessionService {
	if instSessionService != nil {
		return instSessionService
	}

	if kv := store.GetStore(); kv == nil { // 单例
		logger.Errorf("msg kv is not initialized")
		return nil
	} else {
		instSessionService = &SessionService{kv: kv}
		return instSessionService
	}
}

// GetLatestSession 获取访客最新会话（按 session_seq 最大）
func (s *SessionService) GetLatestSession(visitorID, appID string) (*models.Session, error) {
	prefix := fmt.Sprintf("s:%s:%s:", visitorID, appID)
	var latestSession *models.Session

	err := s.kv.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			Prefix:  []byte(prefix),
			Reverse: true, // 从最大 seq 开始
		})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				continue
			}
			var sess models.Session
			if err := json.Unmarshal(val, &sess); err != nil {
				continue
			}
			latestSession = &sess
			return nil // 找到第一个就是最新的
		}
		return fmt.Errorf("session not found for visitor=%s app=%s", visitorID, appID)
	})

	if err != nil {
		logger.Errorf("GetLatestSession failed: %v", err)
		return nil, err
	}
	return latestSession, nil
}

// CreateSession 创建新会话
func (s *SessionService) CreateSession(visitorID, appID string) (*models.Session, error) {
	sessionSeq, err := store.NewSessionSeq() // 需要你实现这个全局序号生成器
	if err != nil {
		logger.Errorf("generate session seq failed: %v", err)
		return nil, err
	}

	sid := models.GetSessionID(visitorID, appID, sessionSeq)
	now := time.Now().Unix()

	session := &models.Session{
		SID:                sid,
		CreatedAt:          now,
		LastVisitorMsgTime: 0,
		LastAgentReplyTime: 0,
		LastAgentReadTime:  0,
		Closed:             false,
		FollowUp:           false,
	}

	data, _ := json.Marshal(session)

	err = s.kv.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(sid), data)
	})
	if err != nil {
		logger.Errorf("save session %s failed: %v", sid, err)
		return nil, err
	}

	return session, nil
}

// GetOrCreateSession 获取或创建会话
func (s *SessionService) GetOrCreateSession(visitorID, appID string) (*models.Session, error) {
	session, err := s.GetLatestSession(visitorID, appID)
	if err == nil {
		// 计算最后活跃时间
		lastActive := session.LastVisitorMsgTime
		if session.LastAgentReplyTime > lastActive {
			lastActive = session.LastAgentReplyTime
		}

		// 如果会话未关闭，但已超时 → 自动关闭并新建
		if !session.Closed && time.Since(time.Unix(lastActive, 0)) > SessionTimeout {
			// 自动关闭旧会话
			session.Close()
			s.SaveSession(session) // 持久化关闭状态

			// 创建新会话
			return s.CreateSession(visitorID, appID)
		}

		// 未超时且未关闭 → 复用
		if !session.Closed {
			return session, nil
		}
	}

	// 无有效会话 → 创建新会话
	return s.CreateSession(visitorID, appID)
}

// SaveSession 保存会话（用于状态更新）
func (s *SessionService) SaveSession(session *models.Session) error {
	if session.SID == "" {
		logger.Errorf("session SID is empty")
		return fmt.Errorf("session SID is empty")
	}
	data, _ := json.Marshal(session)

	err := s.kv.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(session.SID), data)
	})
	if err != nil {
		logger.Errorf("save session %s failed: %v", session.SID, err)
	}
	return err
}

// 获取会话内容
func (s *SessionService) GetSession(sessionID string) (*models.Session, error) {
	var session *models.Session
	if sessionID == "" {
		logger.Errorf("sessionID is empty")
		return nil, fmt.Errorf("sessionID is empty")
	}

	if s.kv == nil {
		logger.Errorf("kv is not initialized")
		return nil, fmt.Errorf("kv is not initialized")
	}

	err := s.kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(sessionID))
		if err != nil {
			logger.Errorf("badger.Get %s failed %v ", sessionID, err)
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			logger.Errorf("badger.Get %s failed %v ", sessionID, err)
			return err
		}
		if err := json.Unmarshal(val, &session); err != nil {
			logger.Errorf("badger.Get %s failed %v ", sessionID, err)
			return err
		}
		return nil
	})

	if err != nil {
		logger.Errorf("GetSession failed: %v", err)
		return nil, err
	}
	return session, nil
}
