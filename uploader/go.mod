module DC-homework-1/uploader

go 1.13

require (
	github.com/gorilla/mux v1.7.4
	github.com/jackc/pgconn v1.6.0 // indirect
	github.com/streadway/amqp v1.0.0
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	lib v0.0.0-00010101000000-000000000000
)

replace lib => ../lib
