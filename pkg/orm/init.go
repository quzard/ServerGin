package orm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewDB(path string) error{
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	return err
}
