package entity

type Product struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
