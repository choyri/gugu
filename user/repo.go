package user

import (
	"context"

	apiv1 "github.com/choyri/gugu/api/v1"
	"github.com/uptrace/bun"
)

type DBRepo struct {
	db bun.IDB
}

var _ Repo = (*DBRepo)(nil)

func NewDBRepo(db bun.IDB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}

func (r *DBRepo) ListUsers(ctx context.Context) ([]*apiv1.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *DBRepo) GetUser(ctx context.Context, name string) (*apiv1.User, error) {
	//TODO implement me
	panic("implement me")
}
