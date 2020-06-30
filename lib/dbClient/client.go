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
			uniqueId TEXT NOT NULL,
			title TEXT NOT NULL,
			category TEXT NOT NULL
		)`); err != nil {
		return Client{}, err
	}
	return
}

func (db *Client) CreateProduct(product Product) (p Product, err error) {
	if product.UniqueID == "" {
		product.HashProduct()
	}
	err = db.conn.QueryRow(
		context.Background(),
		`INSERT INTO products (uniqueId, title, category) VALUES ($1, $2, $3)
			RETURNING id
		`, product.UniqueID, product.Title, product.Category).Scan(&product.ID)
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

func (db *Client) GetProducts(limit, offset int) (products []Product, err error) {
	query := "SELECT id, uniqueId, title, category FROM products ORDER BY id DESC"
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
		if err = rows.Scan(&product.ID, &product.UniqueID, &product.Title, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return
}

func (db *Client) GetProduct(id uint64) (product Product, err error) {
	err = db.conn.QueryRow(
		context.Background(),
		`SELECT id, uniqueId, title, category FROM products WHERE id=$1
		`, id).Scan(&product.ID, &product.UniqueID, &product.Title, &product.Category)
	return
}

func (db *Client) UpdateProduct(product Product) (err error) {
	if product.UniqueID == "" {
		product.HashProduct()
	}
	_, err = db.conn.Exec(
		context.Background(),
		`UPDATE products SET uniqueId=$1, title=$2, category=$3 WHERE id=$4
		`, product.UniqueID, product.Title, product.Category, product.ID)
	return err
}

func (db *Client) CountProducts() (rows int64, err error) {
	err = db.conn.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM products
			`).Scan(&rows)
	return
}

func (db *Client) ImportProductBatch(batch []Product) (err error) {
	productBatch := &pgx.Batch{}
	query := "INSERT INTO products (uniqueId, title, category) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING"
	for _, product := range batch {
		if product.UniqueID == "" {
			product.HashProduct()
		}
		productBatch.Queue(query, product.UniqueID, product.Title, product.Category)
	}
	sendBatch := db.conn.SendBatch(context.Background(), productBatch)
	defer sendBatch.Close()
	_, err = sendBatch.Exec()
	return
}
