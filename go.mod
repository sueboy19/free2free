module free2free

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gorilla/sessions v1.2.1
	github.com/joho/godotenv v1.4.0
	github.com/markbates/goth v1.7.1
	github.com/markbates/goth/providers/facebook v1.7.1
	github.com/markbates/goth/providers/instagram v1.7.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/go-sql-driver/mysql v1.7.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/stretchr/testify v1.8.4
	golang.org/x/oauth2 v0.0.0-20220411215726-087913b9d011
)

// 使用 go mod tidy 來整理相依性