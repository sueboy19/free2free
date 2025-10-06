package main_test

import (
	"fmt"
	"os"
	"testing"

	"free2free/database"
)

// TestMain 設定測試環境
func TestMain(m *testing.M) {
	// Note: Skip actual database initialization for tests 
	// that don't require main package functionality.
	// For tests that need database access, implement
	// appropriate mock or test database setup.
	
	// 執行測試
	code := m.Run()

	// 清理資源
	if database.GlobalDB != nil && database.GlobalDB.Conn != nil {
		sqlDB, err := database.GlobalDB.Conn.DB()
		if err == nil && sqlDB != nil {
			sqlDB.Close()
		}
	}

	os.Exit(code)
}
