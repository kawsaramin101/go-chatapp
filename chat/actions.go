package chat

import (
	db "chatapp/db"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ErrorData struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type CreateChatData struct {
	Usernames []string `json:"usernames"`
}

func CreateChat(msg *Message, c *Client) {
	var data CreateChatData
	err := json.Unmarshal(msg.Data, &data)

	var errorData ErrorData

	if err == nil {
		newChat := db.Chat{SecondaryID: uuid.New().String(), IsPrivateChat: true}

		db.DB.Create(&newChat)

		admin := db.ChatMember{SecondaryID: uuid.New().String(), Role: "admin", UserID: c.dbUser.ID, ChatID: newChat.ID}
		db.DB.Create(&admin)
		db.DB.Model(&newChat).Association("Users").Append(&c.dbUser)
		db.DB.Model(&c.dbUser).Association("Chats").Append(&newChat)

		var newChatMembers []db.ChatMember
		var foundUsers []db.User

		for _, anotherUserUsername := range data.Usernames {
			var anotherUser db.User

			if err := db.DB.Where("username = ?", anotherUserUsername).First(&anotherUser).Error; err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {
					errorData.Action = "ERROR_USER_NOT_FOUND"
					errorData.Message = "User not found"

					jsonData, _ := json.Marshal(errorData)
					c.send <- jsonData
				} else {
					errorData.Action = "ERROR_SERVER_ERROR"
					errorData.Message = "SERVER ERROR"

					jsonData, _ := json.Marshal(errorData)
					c.send <- jsonData
				}
				return

			} else {
				foundUsers = append(foundUsers, anotherUser)

				newChatMember := db.ChatMember{SecondaryID: uuid.New().String(), Role: "user", UserID: anotherUser.ID, ChatID: newChat.ID}

				newChatMembers = append(newChatMembers, newChatMember)

				db.DB.Model(&anotherUser).Association("Chats").Append(&newChat)

			}
			db.DB.Create(&newChatMembers)
			db.DB.Model(&newChat).Association("Users").Append(&foundUsers)

		}
		data := struct {
			Action string `json:"action"`
			Data   struct {
				ChatId          uint   `json:"chatId"`
				ChatSecondaryId string `json:"chatSecondaryId"`
			} `json:"data"`
		}{
			Action: "CHAT_CREATED",
			Data: struct {
				ChatId          uint   `json:"chatId"`
				ChatSecondaryId string `json:"chatSecondaryId"`
			}{
				ChatId:          newChat.ID,
				ChatSecondaryId: newChat.SecondaryID,
			},
		}

		enCodedData, _ := json.Marshal(data)
		c.send <- enCodedData

	} else {
		log.Printf("error unmarshaling JSON: %v", err)
		errorData.Action = "ERROR_INVALID_PAYLOAD"
		errorData.Message = "Usernames not provided"

	}

	enCodedData, _ := json.Marshal(errorData)
	c.send <- enCodedData

}

type CheckIfUserExistData struct {
	Username string `json:"username"`
}

func CheckIfUserExist(msg *Message, c *Client) {

	var data CheckIfUserExistData

	err := json.Unmarshal(msg.Data, &data)

	var errorData ErrorData

	if err != nil {
		var user db.User
		err := db.DB.Where("username = ?", data.Username).First(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				data := map[string]interface{}{
					"action": "CHECK_IF_USER_EXIST",
					"data": map[string]interface{}{
						"username": data.Username,
						"exists":   false,
					},
				}

				encodedData, _ := json.Marshal(data)
				c.send <- encodedData
				return
			} else {
				errorData.Action = "ERROR_SERVER_ERROR"
				errorData.Message = "SERVER ERROR"
			}
		} else {
			data := map[string]interface{}{
				"action": "CHECK_IF_USER_EXIST",
				"data": map[string]interface{}{
					"username": data.Username,
					"exists":   true,
				},
			}

			encodedData, _ := json.Marshal(data)
			c.send <- encodedData

			return
		}

	} else {
		errorData.Action = "ERROR_INVALID_PAYLOAD"
		errorData.Message = "Username not provided"
	}
	encodedData, _ := json.Marshal(errorData)
	c.send <- encodedData
}

type HandleMessageData struct {
	ChatId  float64 `json:"chatId"`
	Message string  `json:"message"`
}

func HandleMessage(msg *Message, c *Client) {
	var data HandleMessageData
	fmt.Println(string(msg.Data))
	err := json.Unmarshal(msg.Data, &data)

	var errorData ErrorData

	if err == nil {
		found := false
		var currentRoom *Room
		for room := range c.rooms {
			if room.dbRoomID == uint(data.ChatId) {
				currentRoom = room
				found = true
				break
			}
		}

		if found {
			data := struct {
				Action string `json:"action"`
				Data   struct {
					ChatId          uint   `json:"chatId"`
					ChatSecondaryId string `json:"chatSecondaryId"`
					Message         string `json:"message"`
					From            string `json:"from"`
				} `json:"data"`
			}{
				Action: "MESSAGE",
				Data: struct {
					ChatId          uint   `json:"chatId"`
					ChatSecondaryId string `json:"chatSecondaryId"`
					Message         string `json:"message"`
					From            string `json:"from"`
				}{
					ChatId:          currentRoom.dbRoomID,
					ChatSecondaryId: currentRoom.dbRoomSecondaryID,
					Message:         data.Message,
					From:            c.dbUser.Username,
				},
			}

			encodedData, err := json.Marshal(data)
			if err != nil {
				errorData.Action = "ERROR_SERVER_ERROR"
				errorData.Message = "Error in JSON encoding"
			} else {
				for client := range currentRoom.clients {
					client.send <- encodedData
				}
				return
			}
		} else {
			errorData.Action = "ERROR_USER_NOT_IN_THE_ROOM"
			errorData.Message = "User is not in the room."

		}
	} else {
		errorData.Action = "ERROR_INVALID_PAYLOAD"
		errorData.Message = "ChatId or message not provided"

	}
	log.Printf("error unmarshaling JSON: %v", err)
	encodedData, _ := json.Marshal(errorData)

	c.send <- encodedData
}
