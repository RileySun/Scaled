module github.com/Scaled/concurrency/errgroup

go 1.23.0

toolchain go1.23.9

require (
	github.com/RileySun/Scaled/utils v0.0.0-20250515131312-dd038c15278a
	golang.org/x/sync v0.13.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace github.com/RileySun/Scaled/utils => ../../utils
