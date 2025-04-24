module github.com/RileySun/monitoring/basic

go 1.23.0

toolchain go1.23.8

require (
	github.com/RileySun/Scaled/utils v0.0.0-20250418122452-b59854248853
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/sync v0.13.0 // indirect
)

replace github.com/RileySun/Scaled/utils => ../../utils
