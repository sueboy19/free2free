package main_test

import (
	"os"
	"testing"
)

// TestMain 設定測試環境
func TestMain(m *testing.M) {
	// 執行測試
	code := m.Run()

	os.Exit(code)
}
