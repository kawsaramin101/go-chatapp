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
	SecondaryID string `gorm:"unique;not null" json:"secondary_id"`
	Username    string `gorm:"not null" json:"username"`
	Password    string `gorm:"not null" json:"-"`

	// Back link to ConnectionRequest
	SentRequests     []ConnectionRequest `gorm:"foreignKey:SendByID"`
	ReceivedRequests []ConnectionRequest `gorm:"foreignKey:SendToID"`

	// Back link to ChatMember
	Chats []Chat `gorm:"many2many:chat_members;"`
}

// Connection model
type ConnectionRequest struct {
	gorm.Model
	SecondaryID string `gorm:"unique;not null" json:"secondaryID"`

	IsAccepted sql.NullBool `gorm:"column:is_accepted" json:"isAccepted"`

	// Foreign keys linking to User model
	SendByID uint `gorm:"not null" json:"-"`
	SendBy   User `gorm:"foreignKey:SendByID" json:"sendBy"`

	SendToID uint `gorm:"not null" json:"-"`
	SendTo   User `gorm:"foreignKey:SendToID" json:"-"`

	// Foreign key linking to Chat model
	ChatID uint `gorm:"not null" json:"-"`
	Chat   Chat `gorm:"foreignKey:ChatID" json:"chat"`
}

type ChatMember struct {
	gorm.Model
	SecondaryID string `gorm:"unique;not null"`
	Role        string `gorm:"type:enum('admin', 'user');default:'user';not null"`

	// Foreign key linking to User model
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`

	// Foreign key linking to Chat model
	ChatID uint `gorm:"not null"`
	Chat   Chat `gorm:"foreignKey:ChatID"`
}

type Chat struct {
	gorm.Model
	SecondaryID string `gorm:"unique;not null" json:"secondary_id"`

	Name string ` json:"name"`

	// If private chat only two people can join, else more than two
	IsPrivateChat bool `gorm:"default:true" json:"is_private_chat"`

	// Back link to ConnectionRequest
	ConnectionRequests []ConnectionRequest `gorm:"foreignKey:ChatID"`

	// Many-to-many relationship with User through ChatMember
	Users []User `gorm:"many2many:chat_members;" json:"users"`
}

func InitializeDB(dataSourceName string) error {

	var err error
	DB, err = gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	err = DB.AutoMigrate(&User{}, &ConnectionRequest{}, &ChatMember{}, &Chat{})
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
