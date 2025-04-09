module github.com/RileySun/Scaled/loadbalancer

go 1.22.4

replace github.com/RileySun/Scaled/utils => ../utils

require github.com/RileySun/Scaled/utils v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)
