package chat

import (
	db "chatapp/db"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateChat(msg *Message, c *Client) {
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

func CheckIfUserExist(msg *Message, c *Client) {
	username, ok := msg.Data["username"].(string)

	errorData := make(map[string]string)
	if ok {
		var user db.User
		err := db.DB.Where("username = ?", username).First(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorData["action"] = "CHECK_IF_USER_EXIST"
				errorData["message"] = "User not found"

				data := map[string]interface{}{
					"action": "CHECK_IF_USER_EXIST",
					"data": map[string]interface{}{
						"username": username,
						"exists":   false,
					},
				}

				encodedData, _ := json.Marshal(data)
				c.send <- encodedData
			} else {
				errorData["action"] = "ERROR_SERVER_ERROR"
				errorData["message"] = "SERVER ERROR"
			}
			return
		} else {
			data := map[string]interface{}{
				"action": "CHECK_IF_USER_EXIST",
				"data": map[string]interface{}{
					"username": username,
					"exists":   true,
				},
			}

			encodedData, _ := json.Marshal(data)
			c.send <- encodedData
		}

	} else {
		errorData["action"] = "ERROR_INVALID_PAYLOAD"
		errorData["message"] = "Username not provided"
	}
	encodedData, _ := json.Marshal(errorData)
	c.send <- encodedData
}

func HandleMessage(msg *Message, c *Client) {
	chatId, ok := msg.Data["chatId"].(float64)
	message, ok := msg.Data["message"].(string)

	errorData := make(map[string]string)

	if ok {
		found := false
		var currentRoom *Room
		for room := range c.rooms {
			if room.dbRoomID == uint(chatId) {
				currentRoom = room
				found = true
				break
			}
		}

		if found {
			data := map[string]interface{}{
				"action": "MESSAGE",
				"data": map[string]interface{}{
					"chatId":          currentRoom.dbRoomID,
					"chatSecondaryId": currentRoom.dbRoomSecondaryID,
					"message":         message,
					"from":            c.dbUser.Username,
				},
			}

			encodedData, err := json.Marshal(data)
			if err != nil {
				errorData["action"] = "ERROR_SERVER_ERROR"
				errorData["message"] = "Error in JSON encoding"
			} else {
				currentRoom.broadcast <- encodedData
				return
			}
		} else {
			errorData["action"] = "ERROR_USER_NOT_IN_THE_ROOM"
			errorData["message"] = "User is not in the room."

		}
	} else {
		errorData["action"] = "ERROR_INVALID_PAYLOAD"
		errorData["message"] = "ChatId or message not provided"

	}

	encodedData, _ := json.Marshal(errorData)

	c.send <- encodedData

}
