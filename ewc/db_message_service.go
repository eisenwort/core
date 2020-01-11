package ewc

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type DbMessageService struct {
	BaseDbService
	updateChatChan chan int64
}

func NewDbMessageService() *DbMessageService {
	srv := new(DbMessageService)
	srv.updateChatChan = make(chan int64, chanSize)

	srv.dbExec(func(db *gorm.DB) {
		db.AutoMigrate(&Message{})
	})

	go srv.listeners()
	return srv
}

func (srv *DbMessageService) Create(msg *Message) (*Message, error) {
	srv.updateChatChan <- msg.ChatID

	var msgError error = nil
	msg.CreatedAt = time.Now()
	//msg.Alghorithm = alghorinthm

	srv.dbExec(func(db *gorm.DB) {
		if err := db.Save(msg).Error; err != nil {
			msgError = errors.New("Ошибка отправки сообщения")
			log.Println("create message error:", err)
			msg = nil
		}
	})

	return msg, msgError
}

func (srv *DbMessageService) Delete(msg *Message) bool {
	srv.updateChatChan <- msg.ChatID
	result := true

	srv.dbExec(func(db *gorm.DB) {
		if err := db.Delete(Message{}, "id = ?", msg.ID).Error; err != nil {
			log.Println("delete message error:", err)
			result = false
		}
	})

	return result
}

// TODO: page
func (srv *DbMessageService) GetByChat(chatID int64) []*Message {
	result := make([]*Message, 0)

	srv.dbExec(func(db *gorm.DB) {
		if err := db.Where(Message{ChatID: chatID}).Order("created_at desc").Find(&result).Error; err != nil {
			log.Println("get message by chat error:", err)
			result = nil
		}
	})

	return result
}

func (srv *DbMessageService) listeners() {
	for {
		select {
		case id := <-srv.updateChatChan:
			srv.dbExec(func(db *gorm.DB) {
				db.Raw("update chat updated_at=now() where chat_id = ?", id)
			})
		}
	}
}
