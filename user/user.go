package user

import (
	"context"

	apiv1 "github.com/choyri/gugu/api/v1"
	"github.com/choyri/gugu/internal/logx"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	apiv1.UnimplementedUserServiceServer
	repo   Repo
	logger *slog.Logger
}

type Repo interface {
	ListUsers(context.Context) ([]*apiv1.User, error)
	GetUser(context.Context, string) (*apiv1.User, error)
}

func NewServer(repo Repo) *Server {
	return &Server{
		repo:   repo,
		logger: logx.NewPkg("user", "server"),
	}
}

func (svr *Server) ListUsers(ctx context.Context, _ *emptypb.Empty) (*apiv1.ListUsersResponse, error) {
	users, err := svr.repo.ListUsers(ctx)
	if err != nil {
		panic("implement me")
	}

	return &apiv1.ListUsersResponse{Users: users}, nil
}

func (svr *Server) GetUser(ctx context.Context, req *apiv1.GetUserRequest) (*apiv1.User, error) {
	panic("implement me")
}

func (svr *Server) CreateUser(ctx context.Context, req *apiv1.CreateUserRequest) (*apiv1.User, error) {
	return &apiv1.User{}, nil
}

func (svr *Server) UpdateUser(ctx context.Context, req *apiv1.UpdateUserRequest) (*apiv1.User, error) {
	panic("implement me")
}
