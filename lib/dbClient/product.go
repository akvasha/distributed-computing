package dbClient

import (
	"crypto/sha256"
	"fmt"
)

type Product struct {
	Title    string `json:"title"`
	ID       uint64 `json:"id"`
	UniqueID string `json:"unique_id"`
	Category string `json:"category"`
}

type encodeFields struct {
	Title    string `json:"title"`
	Category string `json:"category"`
}

func (p *Product) HashProduct() {
	h := sha256.New()
	encode := encodeFields{
		Title:    p.Title,
		Category: p.Category,
	}
	h.Write([]byte(fmt.Sprintf("%v", encode)))
	p.UniqueID = fmt.Sprintf("%x", h.Sum(nil))
}
