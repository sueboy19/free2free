package main

import (
	"os"
	"testing"
)

// TestMain 設定測試環境
func TestMain(m *testing.M) {
	// 在測試開始前執行的程式碼
	// 例如設定測試資料庫連線
	
	// 執行測試
	code := m.Run()
	
	// 在測試結束後執行的程式碼
	// 例如清理測試資料
	
	// 退出測試
	os.Exit(code)
}