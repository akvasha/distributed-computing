package dbClient

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"strings"
	"time"
)

var ErrNotFound = errors.New("Not found")

type Client struct {
	conn *pgx.Conn
}

type User struct {
	Username string
	Password string
	Email    string
}

func InitClient() (dbClient Client, err error) {
	dbClient = Client{}
	dbClient.conn, err = pgx.Connect(
		context.Background(),
		os.Getenv("POSTGRES_ADDRESS"),
	)
	if err != nil {
		return Client{}, err
	}
	if _, err = dbClient.conn.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS users(
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			email TEXT NOT NULL
		)`); err != nil {
		return Client{}, err
	}
	if _, err = dbClient.conn.Exec(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS tokens(
			token TEXT PRIMARY KEY,
			type INTEGER NOT NULL,
			lifetime TIMESTAMP NOT NULL,
			username TEXT NOT NULL
		)`); err != nil {
		return Client{}, err
	}
	return
}

var duplicate = "duplicate key"

func (db *Client) AddUser(user User) (err error) {
	_, err = db.conn.Exec(
		context.Background(),
		`INSERT INTO users (username, password, email) VALUES ($1, $2, $3)
		`, user.Username, user.Password, user.Email)
	if strings.Contains(fmt.Sprintln(err), duplicate) {
		return errors.New("Such user already exists")
	}
	return
}

func (db *Client) GetUser(username string) (user User, err error) {
	if err = db.conn.QueryRow(
		context.Background(),
		`SELECT username, password, email FROM users WHERE username=$1
		`, username).Scan(&user.Username, &user.Password, &user.Email); err == pgx.ErrNoRows {
		return User{}, ErrNotFound
	}
	return
}

type TokenData struct {
	Token    string
	Type     int
	Lifetime time.Time
	Username string
}

func (db *Client) GetToken(token string) (tokenData TokenData, err error) {
	if err = db.conn.QueryRow(
		context.Background(),
		`SELECT token, type, lifetime, username FROM tokens WHERE token=$1
		`, token).Scan(&tokenData.Token, &tokenData.Type,
		&tokenData.Lifetime, &tokenData.Username); err == pgx.ErrNoRows {
		return TokenData{}, ErrNotFound
	}
	return
}

var ErrorDuplicateToken = errors.New("Duplicate token")

func (db *Client) AddToken(token TokenData) (err error) {
	_, err = db.conn.Exec(
		context.Background(),
		`INSERT INTO tokens (token, type, lifetime, username) VALUES ($1, $2, $3, $4)
		`, token.Token, token.Type, token.Lifetime, token.Username)
	if strings.Contains(fmt.Sprintln(err), duplicate) {
		return ErrorDuplicateToken
	}
	return
}
