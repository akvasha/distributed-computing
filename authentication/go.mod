module DC-homework-1/authentication

go 1.13

require (
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/jackc/pgx/v4 v4.6.0
	github.com/streadway/amqp v1.0.0
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/sys v0.0.0-20200610111108-226ff32320da // indirect
	google.golang.org/genproto v0.0.0-20200612171551-7676ae05be11 // indirect
	lib v0.0.0-00010101000000-000000000000
)

replace lib => ../lib
