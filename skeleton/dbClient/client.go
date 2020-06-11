package dbClient

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type Client struct {
	conn *pgx.Conn
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
		`CREATE TABLE IF NOT EXISTS products(
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			category TEXT NOT NULL
		)`); err != nil {
		return Client{}, err
	}
	return
}

func (db *Client) CreateProduct(product Product) (p Product, err error) {
	err = db.conn.QueryRow(
		context.Background(),
		`INSERT INTO products (title, category) VALUES ($1, $2)
			RETURNING id
		`, product.Title, product.Category).Scan(&product.ID)
	p = product
	return
}

func (db *Client) DeleteProduct(id uint64) (err error) {
	_, err = db.conn.Exec(
		context.Background(),
		`DELETE FROM products WHERE id=$1
		`, id)
	return
}

func (db *Client) GetProducts(limit, offset uint64) (products []Product, err error) {
	query := "SELECT id, title, category FROM products ORDER BY id DESC"
	if limit != 0 {
		query += fmt.Sprintf(` LIMIT %d`, limit)
	}
	if offset != 0 {
		query += fmt.Sprintf(` OFFSET %d`, offset)
	}
	rows, err := db.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	products = make([]Product, 0)
	for rows.Next() {
		var product Product
		if err = rows.Scan(&product.ID, &product.Title, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return
}

func (db *Client) GetProduct(id uint64) (product Product, err error) {
	err = db.conn.QueryRow(
		context.Background(),
		`SELECT id, title, category FROM products WHERE id=$1
		`, id).Scan(&product.ID, &product.Title, &product.Category)
	return
}

func (db *Client) UpdateProduct(product Product) (err error) {
	_, err = db.conn.Exec(
		context.Background(),
		`UPDATE products SET title=$1, category=$2 WHERE id=$3
		`, product.Title, product.Category, product.ID)
	return err
}
