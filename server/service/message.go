package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"

	"kefu-server/models"
	"kefu-server/store"
	"kefu-server/utils/logger"
)

// 1. 消息存储
// m:{visitor_id}:{app_id}:{session_seq}:{msg_seq}
// m:alice:shop123:1234:23445
// │   │     │       │      └─ msg seq
// │   │     │       └─ session seq
// │   │     └─ app_id
// │   └─ visitor_id

type MessageService struct {
	kv *badger.DB
}

var (
	instMessage *MessageService
)

func GetMsgService() *MessageService {
	if instMessage != nil {
		return instMessage
	}

	if kv := store.GetStore(); kv == nil { // 单例
		return nil
	} else {
		instMessage = &MessageService{kv: kv}
		return instMessage
	}
}

// SaveMessage 保存消息
func (m *MessageService) SaveMessage(visitorID, appID string, sessionSeq uint32, msg *models.Message) (string, error) {
	msgSeq, err := store.NewMsgSeq()
	if err != nil {
		logger.Errorf("getNextID %v ", err)
		return "", err
	}
	msgID := models.GetMessageID(visitorID, appID, sessionSeq, msgSeq)
	msg.MsgID = msgID
	data, _ := json.Marshal(msg)

	// 设置 TTL：30 天
	ttl := 30 * 24 * time.Hour

	err = m.kv.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(msgID), data).WithTTL(ttl)
		return txn.SetEntry(entry)
	})
	if err != nil {
		logger.Errorf("badger.Update %v ", err)
		return "", err
	}
	return msgID, err
}
func (m *MessageService) GetMessage(msgID string) (*models.Message, error) {
	var msg models.Message
	err := m.kv.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(msgID))
		if err != nil {
			logger.Errorf("badger.Get %s failed %v ", msgID, err)
			return err
		}
		val, _ := item.ValueCopy(nil)
		return json.Unmarshal(val, &msg)
	})
	if err != nil {
		logger.Errorf("read msg %s failed %v ", msgID, err)
		return nil, err
	}
	return &msg, err
}

// GetMessagesBySession 获取某会话的消息列表（按时间正序）
func (m *MessageService) GetMessagesBySession(sessionID string, limit int) ([]*models.Message, error) {
	if limit <= 0 || limit > 100 { // 防止滥用
		limit = 50
	}

	// 1. 解析 sessionID: s:{visitor}:{app}:{session_seq}
	parts := strings.Split(sessionID, ":")
	if len(parts) != 4 || parts[0] != "s" {
		return nil, fmt.Errorf("invalid sessionID format: %s", sessionID)
	}
	visitorID := parts[1]
	appID := parts[2]
	sessionSeqStr := parts[3]

	// 2. 构造消息 key 前缀: m:{visitor}:{app}:{session_seq}:
	msgPrefix := fmt.Sprintf("m:%s:%s:%s:", visitorID, appID, sessionSeqStr)

	var msgs []*models.Message

	err := m.kv.View(func(txn *badger.Txn) error {
		opts := badger.IteratorOptions{
			Prefix:  []byte(msgPrefix),
			Reverse: true, // 从最新消息开始扫
		}
		it := txn.NewIterator(opts)
		defer it.Close()

		count := 0
		for it.Rewind(); it.Valid() && count < limit; it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				continue
			}
			var msg models.Message
			if err := json.Unmarshal(val, &msg); err != nil {
				continue
			}
			msgs = append(msgs, &msg)
			count++
		}
		return nil
	})

	if err != nil {
		logger.Errorf("GetMessagesBySession failed: %v", err)
		return nil, err
	}

	// 3. 反转，使结果为时间正序（旧 → 新）
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nil
}

// Close 关闭 DB
func (m *MessageService) Close() error {
	m.kv.Close()
	m.kv = nil
	return nil
}
