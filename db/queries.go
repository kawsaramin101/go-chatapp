package db

import "github.com/google/uuid"

func createNewChat(name string, isPrivateChat bool) *Chat {
	newChat := Chat{SecondaryID: uuid.New().String(), Name: name, IsPrivateChat: isPrivateChat}

	DB.Create(&newChat)

	return &newChat

}
