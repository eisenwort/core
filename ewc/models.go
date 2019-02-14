package ewc

//go:generate go-mobile-collection $GOFILE
import (
	"io"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name          string `gorm:"column:name" sql:"index"`
	Password      string `gorm:"column:password"`
	ResetPassword string `gorm:"column:reset_password"`
}

func (User) TableName() string {
	return "users"
}

// @collection-wrapper
type Message struct {
	gorm.Model
	Text       string `gorm:"column:text"`
	ImageUrl   string `gorm:"column:image_url"`
	Alghorithm string `gorm:"column:alghorithm"`
	SenderId   uint   `gorm:"column:password"`
	ReceiverId uint   `gorm:"column:reset_password"`
}

func (Message) TableName() string {
	return "messages"
}

// @collection-wrapper
type Friend struct {
	gorm.Model
	UserId   uint `gorm:"column:user_id" sql:"index"`
	FriendId uint `gorm:"column:friend_id"`
	User     User
}

func (Friend) TableName() string {
	return "friends"
}

type Settings struct {
	gorm.Model
	UserId int64 `gorm:"column:user_id" json:"user_id"`
}

func (Settings) TableName() string {
	return "friends"
}

type ApiRequest struct {
	Method     string
	RequestUrl string
	Token      string
	Body       io.Reader
}
