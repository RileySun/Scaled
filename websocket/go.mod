module github.com/RileySun/Scaled/websocket

go 1.22.4

require (
	github.com/RileySun/Scaled/utils v0.0.0-20250417102003-6abf75280cfc
	github.com/godbus/dbus/v5 v5.1.0
	github.com/gorilla/websocket v1.5.3
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace github.com/RileySun/Scaled/utils => ../utils
