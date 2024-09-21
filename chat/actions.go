package chat

import (
	db "chatapp/db"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateChat(msg Message, c *Client) {
	if anotherUserUsername, ok := msg.Data["username"].(string); ok {
		var anotherUser db.User

		if err := db.DB.Where("username = ?", anotherUserUsername).First(&anotherUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response := map[string]string{
					"action":  "ERROR_USER_NOT_FOUND",
					"message": "User not found",
				}
				jsonData, _ := json.Marshal(response)

				c.send <- jsonData
			} else {
				response := map[string]string{
					"action":  "ERROR_SERVER_ERROR",
					"message": "SERVER ERROR",
				}
				jsonData, _ := json.Marshal(response)

				c.send <- jsonData
			}
		} else {
			newChat := db.Chat{SecondaryID: uuid.New().String(), IsPrivateChat: true}

			db.DB.Create(&newChat)

			chatMembers := []db.ChatMember{
				{SecondaryID: uuid.New().String(), Role: "admin", UserID: c.dbUser.ID, ChatID: newChat.ID},
				{SecondaryID: uuid.New().String(), Role: "user", UserID: anotherUser.ID, ChatID: newChat.ID},
			}

			db.DB.Create(&chatMembers)

			db.DB.Model(&newChat).Association("Users").Append(&c.dbUser, &anotherUser)

			db.DB.Model(&c.dbUser).Association("Chats").Append(&newChat)
			db.DB.Model(&anotherUser).Association("Chats").Append(&newChat)

			data := map[string]interface{}{
				"action": "CHAT_CREATED",
				"data": map[string]interface{}{
					"chatId":          newChat.ID,
					"chatSecondaryId": newChat.SecondaryID,
				},
			}

			enCodedData, _ := json.Marshal(data)

			c.send <- enCodedData
		}
	} else {
		data := map[string]interface{}{
			"action":  "ERROR_INVALID_PAYLOAD",
			"message": "Username not provided",
		}

		enCodedData, _ := json.Marshal(data)

		c.send <- enCodedData
	}

}
