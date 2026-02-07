package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"kefu-server/utils/logger"
)

type logWriter struct{}

var (
	DB *gorm.DB
)

func (w logWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	// If there are multiple lines, remove the first line
	if idx := strings.Index(msg, "\n"); idx != -1 {
		msg = msg[idx+1:]
	}
	// Remove all newlines
	msg = strings.ReplaceAll(msg, "\n", " ")
	// Log as Info level
	logger.Infof("[SQL] %s", msg)
	return len(p), nil
}

// InitDB initializes database connection
func InitDB(dbPath string) (*gorm.DB, error) {
	// Ensure data directory exists
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			logger.Errorf("Failed to create data directory: %v", err)
			return nil, err
		}
	}
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: glogger.New(log.New(logWriter{}, "", 0), glogger.Config{LogLevel: glogger.Info}),
	})
	if err != nil {
		logger.Errorf("Failed to open database: %v", err)
		return nil, err
	}

	return DB, nil
}
