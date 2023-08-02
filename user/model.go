package user

import (
	"time"
)

type User struct {
	ID          int32  `bun:",pk,autoincrement"`
	Name        string `bun:",skipupdate" validate:"required_if=ID 0"`
	DisplayName *string
	Password    *string
	CreateTime  time.Time `bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
	UpdateTime  time.Time `bun:",nullzero,notnull,default:current_timestamp,skipupdate"`
}
