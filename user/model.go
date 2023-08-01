package user

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
}
