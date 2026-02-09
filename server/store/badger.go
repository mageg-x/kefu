package store

import (
	"fmt"
	"kefu-server/utils/logger"
	"strconv"

	"github.com/dgraph-io/badger/v4"
)

var (
	KV *badger.DB
)

func GetStore() *badger.DB { return KV }

func InitStore(dbPath string) (*badger.DB, error) {
	if KV != nil { // 单例
		return KV, nil
	}

	opts := badger.DefaultOptions(dbPath).
		WithLogger(nil).                // 禁用 badger 自带日志（用你的 logger）
		WithValueLogFileSize(128 << 20) // 128MB

	db, err := badger.Open(opts)
	if err != nil {
		logger.Errorf("failed to open badger: %v", err)
		return nil, fmt.Errorf("failed to open badger: %w", err)
	}
	KV = db
	return KV, nil
}

// GetNextID 原子获取下一个 ID（线程安全）
func GetNextID(counterKey string) (uint32, error) {
	if KV == nil {
		return 0, fmt.Errorf("store not initialized")
	}

	var id uint32
	err := KV.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(counterKey))
		if err == badger.ErrKeyNotFound {
			id = 1
		} else if err != nil {
			logger.Errorf("badger.Get %v", err)
			return err
		} else {
			val, _ := item.ValueCopy(nil)
			current, _ := strconv.ParseUint(string(val), 10, 64)
			id = uint32(current + 1)
		}
		return txn.Set([]byte(counterKey), []byte(strconv.FormatUint(uint64(id), 10)))
	})
	return id, err
}

func NewSessionSeq() (uint32, error) {
	seq, err := GetNextID("counter:session")
	if err != nil {
		logger.Errorf("getNextID for session failed %v ", err)
		return 0, err
	}
	return seq, nil
}

func NewMsgSeq() (uint32, error) {
	seq, err := GetNextID("counter:message")
	if err != nil {
		logger.Errorf("getNextID for message failed %v ", err)
		return 0, err
	}
	return seq, nil
}
