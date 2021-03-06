package ewc

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

type DbUserService struct {
}

func NewDbUserService() *DbUserService {
	srv := new(DbUserService)

	dbExec(func(db *gorm.DB) {
		db.AutoMigrate(&User{})
		db.AutoMigrate(&UserData{})
		db.AutoMigrate(&Friend{})
	})

	return srv
}

func (srv *DbUserService) Create(login, password, passwordForReset string) (*User, error) {
	existingUser := new(User)
	var user = new(User)
	var userError error = nil

	dbExec(func(db *gorm.DB) {
		if err := db.Where(User{Login: login}).First(existingUser); err == nil {
			userError = errors.New("Логин уже занят")
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hashedResetPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordForReset), bcrypt.DefaultCost)
		user := &User{
			Login:         login,
			Password:      string(hashedPassword),
			ResetPassword: string(hashedResetPassword),
		}
		if err := db.Save(user).Error; err != nil {
			log.Println("create user error:", err)
			user = nil
			userError = errors.New("Произошла неизвестная ошибка")
		}
	})

	return user, userError
}

func (srv *DbUserService) Update(user *User) *User {
	dbExec(func(db *gorm.DB) {
		if err := db.Save(user).Error; err != nil {
			log.Println("update user error:", err)
		}
	})

	return user
}

func (srv *DbUserService) Save(user *User) *User {
	dbExec(func(db *gorm.DB) {
		if err := db.Save(user).Error; err != nil {
			log.Println("save user error:", err)
		}
	})

	return user
}

func (srv *DbUserService) SaveUserData(data UserData) {
	dbExec(func(db *gorm.DB) {
		db.Delete(&UserData{})

		if err := db.Save(&data).Error; err != nil {
			log.Println("save user data error:", err)
		}
	})
}

func (srv *DbUserService) Migrate() {
	dbExec(func(db *gorm.DB) {
		db.AutoMigrate(User{})
		db.AutoMigrate(UserData{})
		db.AutoMigrate(Friend{})
	})
}

func (srv *DbUserService) Login(login, password string) *User {
	user := new(User)

	dbExec(func(db *gorm.DB) {
		if err := db.Where(User{Login: login}).First(&user).Error; err != nil {
			log.Println("get user error:", err)
		}
	})

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		return user
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.ResetPassword), []byte(password)); err == nil {
		dbExec(func(db *gorm.DB) {
			db.Delete(&Message{}, "user_id = ?", user.ID)
			db.Delete(&ChatUser{}, "user_id = ?", user.ID)
			db.Delete(&Chat{}, "owner_id = ?", user.ID)

			user.Reseted = true
			if err := db.Save(user).Error; err != nil {
				log.Println("set user reseted error")
			}
		})
		return user
	}

	return nil
}

func (srv *DbUserService) IsLogin() User {
	user := User{}

	dbExec(func(db *gorm.DB) {
		db.Model(User{}).First(&user)
	})

	return user
}

func (srv *DbUserService) GetFriends(userID int64) []User {
	users := make([]User, 0)

	dbExec(func(db *gorm.DB) {
		db.Model(User{}).
			Where("id in (select friend_id from friends where user_id = ?)", userID).
			Find(&users)
	})

	return users
}

func (srv *DbUserService) SaveFriends(userId int64, users []User) {
	dbExec(func(db *gorm.DB) {
		db.Where(Friend{UserID: userId}).Delete(&Friend{})

		for _, user := range users {
			friend := Friend{
				UserID:   userId,
				FriendID: user.ID,
			}
			db.Save(&friend)
		}
	})
}

func (srv *DbUserService) Get(id int64) User {
	user := User{}

	dbExec(func(db *gorm.DB) {
		db.First(&user, id)
	})

	return user
}

func (srv *DbUserService) GetUserData() UserData {
	userData := UserData{}

	dbExec(func(db *gorm.DB) {
		db.First(&userData)
	})

	return userData
}

func (srv *DbUserService) GetByLogin(login string) User {
	user := User{}

	dbExec(func(db *gorm.DB) {
		db.Where(&User{Login: login}).First(&user)
	})

	return user
}

func (srv *DbUserService) Search(searchString string) []*User {
	users := make([]*User, 0)

	dbExec(func(db *gorm.DB) {
		db.Where("login ilike ?", searchString).Find(&users)
	})

	return users
}

func (srv *DbUserService) AddFriend(userId int64, friendID int64) User {
	user := User{}

	dbExec(func(db *gorm.DB) {
		friendItem := Friend{
			UserID:   userId,
			FriendID: friendID,
		}
		if err := db.Save(&friendItem).Error; err != nil {
			log.Println("save new friend error:", err)
			return
		}
		if err := db.First(&user, friendID).Error; err != nil {
			log.Println("get user for added friend error:", err)
		}
	})

	return user
}

func (srv *DbUserService) DeleteFriend(userID int64, friendID int64) bool {
	success := true

	dbExec(func(db *gorm.DB) {
		if err := db.Where(Friend{UserID: userID, FriendID: friendID}).Delete(Friend{}).Error; err != nil {
			log.Printf("delete friend %d for user %d, error: %s", friendID, userID, err.Error())
			success = false
		}

		chats := []Chat{}
		chatIds := []int64{}
		query := `
			select * from chats c
				join chats_users cu on c.id = cu.chat_id
			where
				c.owner_id = ?
				and c.personal = true
				and cu.user_id in (?)
			ORDER BY c.id;
		`
		// выбрать все персональные чаты где собеседником явлеяется  этот контак
		db.Raw(query, userID, friendID).Scan(&chats)

		for _, chat := range chats {
			chatIds = append(chatIds, chat.ID)
		}

		// выбрать все персональные чаты контакта где собеседником является юзер
		db.Raw(query, friendID, userID).Scan(&chats)

		for _, chat := range chats {
			chatIds = append(chatIds, chat.ID)
		}

		// удалить чаты и связи
		db.Where("id in (?)", chatIds).Delete(&Chat{})
		db.Where("chat_id in (?)", chatIds).Delete(&ChatUser{})
	})

	return success
}

func (srv *DbUserService) Delete(id int64) {
	dbExec(func(db *gorm.DB) {
		if err := db.Delete(Chat{}, "owner_id = ?", id).Error; err != nil {
			log.Println("delete user chats error:", err)
		}
		if err := db.Delete(ChatUser{}, "user_id = ?", id).Error; err != nil {
			log.Println("delete chat users error:", err)
		}
		if err := db.Delete(Message{}, "sender_id = ? and receiver_id = ?", id).Error; err != nil {
			log.Println("delete user messages error:", err)
		}
	})
}
