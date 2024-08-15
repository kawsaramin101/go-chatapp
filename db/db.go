package db

import (
	"database/sql"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	SecondaryID            string        `gorm:"unique;not null"`
	Username               string        `gorm:"not null"`
	Password               string        `gorm:"not null"`
	ConnectionRequestsFrom []*Connection `gorm:"foreignKey:SendBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ConnectionRequestsTo   []*Connection `gorm:"foreignKey:SendTo;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Connection model
type Connection struct {
	gorm.Model
	SecondaryID string       `gorm:"unique;not null"`
	SendBy      uint         `gorm:"not null"`
	SendTo      uint         `gorm:"not null"`
	IsAccepted  sql.NullBool `gorm:"column:is_accepted"`
}

func InitializeDB(dataSourceName string) error {

	var err error
	DB, err = gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(&User{}, &Connection{})
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal(err)
	}
}
