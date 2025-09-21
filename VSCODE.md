# Visual Studio Code 開發設定

## 使用 Dev Container (推薦)

本專案包含 Dev Container 設定，可以讓你在容器化環境中開發，確保環境一致性。

### 設定步驟
1. 安裝 Docker Desktop
2. 安裝 VS Code 的 Dev Container 擴充套件
3. 在 VS Code 中開啟專案資料夾
4. 按下 `Ctrl+Shift+P`，輸入 "Dev Containers: Reopen in Container"

### 環境變數設定
Dev Container 設定中已包含環境變數的範本，請在 `.devcontainer/devcontainer.json` 中設定實際的值。

## 使用 launch.json 執行和除錯

專案包含 VS Code 的 launch.json 設定，可以直接在 VS Code 中執行和除錯應用程式。

### 執行步驟
1. 開啟 VS Code
2. 使用 Docker Compose 啟動資料庫：
   ```bash
   docker-compose up -d
   ```
3. 等待資料庫啟動完成
4. 按下 `Ctrl+Shift+D` 開啟 Run and Debug 面板
5. 選擇 "Launch Package" 設定
6. 按下 "Start Debugging" 按鈕 (綠色箭頭) 或按下 F5

### 環境變數設定
launch.json 設定中已包含環境變數的範本，請在 `.vscode/launch.json` 中設定實際的值。

## 反向代理設定 (devtunnel)

為了方便進行 Facebook 和 Instagram 的 OAuth 登入測試，建議使用 devtunnel 建立反向代理。

### 安裝 devtunnel
1. 安裝 Visual Studio (2022 或更新版本) 或 Visual Studio Code
2. devtunnel 工具會隨 Visual Studio 或 VS Code 一起安裝
3. 如果僅需要命令列工具，可以從 https://aka.ms/devtunnel/cli 安裝

### 使用步驟
1. 啟動應用程式: `go run main.go`
2. 在終端機執行: `devtunnel host 8080`
3. devtunnel 會提供一個公開 URL (例如: https://abcd1234.devtunnel.azure.com)
4. 將此 URL 設定為環境變數 BASE_URL 和 OAuth 應用程式的回調 URL

### 環境變數更新
當使用 devtunnel 時，需要更新以下環境變數：
- BASE_URL: 設定為 devtunnel 提供的 URL (例如: https://abcd1234.devtunnel.azure.com)
- FACEBOOK_KEY, FACEBOOK_SECRET: 在 Facebook 開發者中心設定應用程式的 OAuth 回調 URL
- INSTAGRAM_KEY, INSTAGRAM_SECRET: 在 Instagram 開發者中心設定應用程式的 OAuth 回調 URL