module github.com/RileySun/Scaled/login

go 1.23.0

toolchain go1.23.8

require (
	github.com/RileySun/Scaled/utils v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/julienschmidt/httprouter v1.3.0
	golang.org/x/crypto v0.37.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/sync v0.13.0 // indirect
)

replace github.com/RileySun/Scaled/utils => ../utils
