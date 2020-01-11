package ewc

//go:generate go-mobile-collection $GOFILE
import (
	"crypto/rsa"
	"hash"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID            int64  `gorm:"primary_key" json:"id,omitempty"`
	Login         string `gorm:"column:login" sql:"index" json:"login,omitempty"`
	Password      string `gorm:"column:password"`
	ResetPassword string `gorm:"column:reset_password"`
	PublicKey     string `gorm:"column:public_key"`
	PrivateKey    string `gorm:"column:private_key"`
	Reseted       bool   `gorm:"column:reseted"  json:"reseted,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (v *User) Equal(rhs *User) bool {
	return v.Login == rhs.Login
}

type UserData struct {
	RefreshToken string `gorm:"column:refresh_token"`
}

// @collection-wrapper
type Message struct {
	ID         int64     `gorm:"primary_key" json:"id,omitempty"`
	UserID     int64     `gorm:"column:user_id" json:"sender_id,omitempty"`
	ChatID     int64     `gorm:"column:chat_id" json:"chat_id,omitempty"`
	Text       string    `gorm:"column:text" json:"text,omitempty"`
	ImageURL   string    `gorm:"column:image_url" json:"image_url,omitempty"`
	FileURL    string    `gorm:"column:file_url" json:"file_url,omitempty"`
	Alghorithm string    `gorm:"column:alghorithm" json:"alghorithm,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	ExpiredAt  time.Time `gorm:"column:expired_at" json:"expired_at,omitempty"`
	User       User
	Chat       Chat
}

func (Message) TableName() string {
	return "messages"
}

func (v *Message) Equal(rhs *Message) bool {
	return v.ID == rhs.ID
}

// @collection-wrapper
type Friend struct {
	ID       int64 `gorm:"primary_key" json:"id,omitempty"`
	UserID   int64 `gorm:"column:user_id" sql:"index" json:"user_id,omitempty"`
	FriendID int64 `gorm:"column:friend_id" json:"friend_id,omitempty"`
	User     User  `json:"user,omitempty"`
}

func (Friend) TableName() string {
	return "friends"
}

func (v *Friend) Equal(rhs *Friend) bool {
	return v.UserID == rhs.UserID
}

type Settings struct {
	ID     int64 `gorm:"primary_key" json:"id,omitempty"`
	UserId int64 `gorm:"column:user_id" json:"user_id"`
}

func (Settings) TableName() string {
	return "friends"
}

// @collection-wrapper
type Chat struct {
	ID             int64     `gorm:"primary_key" json:"id,omitempty"`
	OwnerID        int64     `gorm:"column:owner_id" json:"owner_id,omitempty"`
	UnreadMessages int       `gorm:"column:unread_messages" json:"unread_messages,omitempty"`
	Name           string    `gorm:"column:user_id" json:"name,omitempty"`
	Personal       bool      `gorm:"column:personal" json:"personal,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	Users          []User    `json:"users,omitempty"`
	Messages       []Message `json:"messages,omitempty"`
}

func (Chat) TableName() string {
	return "chats"
}

func (v *Chat) Equal(rhs *Chat) bool {
	return v.ID == rhs.ID
}

type ChatUser struct {
	ID     int64 `gorm:"primary_key" json:"id,omitempty"`
	ChatID int64 `gorm:"column:chat_id"`
	UserID int64 `gorm:"column:user_id"`
	Chat   Chat
	User   User
}

func (ChatUser) TableName() string {
	return "chats_users"
}

type ApiRequest struct {
	Method     string
	RequestUrl string
	Token      string
	Body       []byte
}

type SocketMessage struct {
	Key  string      `json:"key,omitempty"`
	Body interface{} `json:"body,omitempty"`
}

type EncryptData struct {
	Message   []byte
	Label     []byte
	PublicKey *rsa.PublicKey
	Hash      hash.Hash
}

type JwtClaims struct {
	*jwt.MapClaims
	Id int64
}

type TokenData struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type SetupData struct {
	DbPath           string
	DbDriver         string
	ConnectionString string
}
