package model

type Todo struct {
	BaseModel
	Title     string `bun:",notnull"`
	Completed bool   `bun:",notnull"`
}
