package ewc

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type DbChatService struct {
	BaseDbService
}

func NewDbChatService() *DbChatService {
	srv := new(DbChatService)

	srv.dbExec(func(db *gorm.DB) {
		db.AutoMigrate(&Chat{})
		db.AutoMigrate(&ChatUser{})
	})

	return srv
}

func (srv *DbChatService) Create(chat *Chat) (*Chat, error) {
	existingChat := Chat{}
	var chatError error = nil
	chat.CreatedAt = time.Now()

	srv.dbExec(func(db *gorm.DB) {
		user := &User{}

		if err := db.First(user, chat.OwnerID).Error; err != nil {
			chatError = errors.New("Произошла неизвестная ошибка")
			return
		}
		if user.Reseted {
			chatError = errors.New("Произошла неизвестная ошибка")
			return
		}
		if err := db.Where(Chat{Name: chat.Name, OwnerID: chat.OwnerID}).First(&existingChat); err == nil {
			chat = &existingChat
			return
		}
		if err := db.Save(&chat).Error; err != nil {
			chatError = errors.New("Произошла неизвестная ошибка")
		}
		for _, user := range chat.Users {
			chatUser := ChatUser{
				ChatID: chat.ID,
				UserID: user.ID,
			}
			if err := db.Save(&chatUser).Error; err != nil {
				chatError = errors.New("Произошла неизвестная ошибка")
			}
		}
	})

	return chat, chatError
}

func (srv *DbChatService) Update(chat *Chat) (*Chat, error) {
	var chatError error = nil
	chat.UpdatedAt = time.Now()

	srv.dbExec(func(db *gorm.DB) {
		if err := db.Save(chat).Error; err != nil {
			log.Println("update user error:", err)
			chatError = errors.New("не удалось обновить чат")
		}
	})

	return chat, chatError
}

func (srv *DbChatService) Get(id int64, include []string) (*Chat, error) {
	var chat = new(Chat)
	var chatError error = nil

	srv.dbExec(func(db *gorm.DB) {
		query := db.Order("updated_at desc")

		if Contains(include, "messages") {
			query = query.Preload("Messages")
		}
		if err := query.First(chat, id).Error; err != nil {
			chatError = errors.New("Невозможно получить чат")
		}
		if chatError == nil && Contains(include, "users") {
			chatUsers := make([]ChatUser, 0)
			chat.Users = make([]User, 0)

			if err := db.Preload("User").Where(ChatUser{ChatID: chat.ID}).Find(&chatUsers).Error; err != nil {
				return
			}
			for _, rel := range chatUsers {
				chat.Users = append(chat.Users, rel.User)
			}
		}
	})

	return chat, chatError
}

func (srv *DbChatService) GetForUser(ownerID int64) ([]*Chat, error) {
	var chats = make([]*Chat, 0)
	var chatError error = nil

	srv.dbExec(func(db *gorm.DB) {
		db = db.Debug()
		if err := db.Where(Chat{OwnerID: ownerID}).Find(&chats).Error; err != nil {
			chatError = errors.New("Невозможно получить чат")
		}
		for _, chat := range chats {
			if !chat.Personal {
				continue
			}

			chatUsers := make([]ChatUser, 0)

			if err := db.Preload("User").Where(ChatUser{ChatID: chat.ID}).Find(&chatUsers).Error; err != nil {
				continue
			}
			for _, rel := range chatUsers {
				if rel.User.ID != userID {
					chat.Name = rel.User.Login
				}
			}
		}
	})

	return chats, chatError
}

func (srv *DbChatService) Delete(chat *Chat) {
	srv.dbExec(func(db *gorm.DB) {
		if err := db.Delete(Chat{}, "id = ?", chat.ID).Error; err != nil {
			log.Println("delete chat error:", err)
		}
		if err := db.Delete(ChatUser{}, "chat_id = ?", chat.ID).Error; err != nil {
			log.Println("delete chat users error:", err)
		}
		if err := db.Delete(Message{}, "chat_id = ?", chat.ID).Error; err != nil {
			log.Println("delete chat messages error:", err)
		}
	})
}

func (srv *DbChatService) Exit(chat *Chat) {
	srv.dbExec(func(db *gorm.DB) {
		if err := db.Delete(ChatUser{}, "chat_id = ? and user_id = ?", chat.ID, userID).Error; err != nil {
			log.Println("delete chat users error:", err)
		}
		if err := db.Delete(Message{}, "chat_id = ? and user_id = ?", chat.ID, userID).Error; err != nil {
			log.Println("delete chat messages error:", err)
		}
	})
}

func (srv *DbChatService) Clean(chat *Chat) bool {
	result := true

	srv.dbExec(func(db *gorm.DB) {
		if err := db.Where(&Message{ChatID: chat.ID}).Delete(Message{}).Error; err != nil {
			log.Println("clean chat error:", err)
			result = false
		}
	})

	return result
}
