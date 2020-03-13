package dbClient

type Product struct {
	Title    string `json:"title"`
	ID       uint64 `json:"id"`
	Category string `json:"category"`
}
