package db

import (
	"em-task/cfg"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDB(config *cfg.Config) *Db {
	db, err := gorm.Open(postgres.Open(config.Db.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
