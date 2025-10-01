# 多階段建置：第一階段編譯 Go 應用
FROM golang:1.25-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 以快取依賴
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源碼
COPY . .

# 建置應用，使用靜態連結以便在 alpine 中運行
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o free2free-staging .

# 第二階段：運行階段，使用輕量 alpine 映像
FROM alpine:latest

# 安裝 ca-certificates 以支援 HTTPS
RUN apk --no-cache add ca-certificates

# 設定工作目錄
WORKDIR /root/

# 從 builder 階段複製建置的二進位檔
COPY --from=builder /app/free2free-staging .

# 安裝 netcat 以等待 DB
RUN apk add --no-cache netcat-openbsd

# 暴露端口
EXPOSE 8080

# 等待 DB 就緒後運行應用
CMD ["sh", "-c", "until nc -z mariadb 3306; do echo 'waiting for database...'; sleep 2; done; ./free2free-staging"]