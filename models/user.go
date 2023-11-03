package model

type User struct {
	BaseModel
	Email    string `bun:",unique,notnull"`
	Name     string `bun:",notnull"`
	Password string `bun:",notnull"`
}
