package entities

import (
	"time"
)

type (
	ChatID  int64
	UserID  int64
	UserIds []UserID // нужно подумать как этот тип можно использовать
	ChatIds []ChatID // нужно подумать как этот тип можно использовать

	User struct {
		UserID       UserID
		LastTime     time.Time
		CountOfChats int64
	}

	Chat struct {
		ChatID       ChatID
		ChatType     uint
		CountOfUsers int64
	}

	ChatList map[ChatID][]User
	UserList map[UserID][]Chat
)
