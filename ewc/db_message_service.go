package ewc

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type DbMessageService struct {
	updateChatChan chan int64
}

func NewDbMessageService() *DbMessageService {
	srv := new(DbMessageService)
	srv.updateChatChan = make(chan int64, chanSize)

	dbExec(func(db *gorm.DB) {
		db.AutoMigrate(&Message{})
	})

	go srv.listeners()
	return srv
}

func (srv *DbMessageService) Create(msg Message) (Message, error) {
	srv.updateChatChan <- msg.ChatID

	var msgError error = nil
	msg.CreatedAt = time.Now()
	//msg.Alghorithm = alghorinthm

	dbExec(func(db *gorm.DB) {
		if err := db.Save(&msg).Error; err != nil {
			msgError = errors.New("Ошибка отправки сообщения")
			log.Println("create message error:", err)
		}
	})

	return msg, msgError
}

func (srv *DbMessageService) Save(msg Message) {
	dbExec(func(db *gorm.DB) {
		if err := db.Save(&msg).Error; err != nil {
			log.Println("create message error:", err)
		}
	})
}

func (srv *DbMessageService) Delete(msg Message) bool {
	srv.updateChatChan <- msg.ChatID
	result := true

	dbExec(func(db *gorm.DB) {
		if err := db.Where("id = ?", msg.ID).Delete(Message{}).Error; err != nil {
			log.Println("delete message error:", err)
			result = false
		}
	})

	return result
}

func (srv *DbMessageService) GetByChat(chatID int64, page int) []Message {
	result := make([]Message, 0)
	srv.updateChatChan <- chatID

	dbExec(func(db *gorm.DB) {
		query := db.Where(Message{ChatID: chatID}).Order("created_at desc")

		if page != 0 {
			offset := pageLimit * (page - 1)
			query = query.Offset(offset).Limit(pageLimit)
		}
		if err := query.Find(&result).Error; err != nil {
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
			dbExec(func(db *gorm.DB) {
				db.Raw("update chat updated_at=now() where chat_id = ?", id)
			})
		}
	}
}
