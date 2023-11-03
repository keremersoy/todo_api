package app

import "github.com/uptrace/bun"

var DB *bun.DB

func Initialise(newDB *bun.DB) {
	DB = newDB
}
