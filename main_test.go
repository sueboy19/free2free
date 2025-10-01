package main

import (
	"fmt"
	"os"
	"testing"
)

// TestMain 設定測試環境
func TestMain(m *testing.M) {
	// 初始化資料庫連接
	if err := InitDB(); err != nil {
		fmt.Printf("初始化資料庫失敗: %v\n", err)
		os.Exit(1)
	}

	// 執行測試
	code := m.Run()

	// 清理資源
	sqlDB, err := GetDB().DB()
	if err == nil && sqlDB != nil {
		sqlDB.Close()
	}

	os.Exit(code)
}
