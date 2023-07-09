package user

import (
	apiv1 "github.com/choyri/gugu/api/v1"
	"github.com/choyri/gugu/internal/logx"
)

type Server struct {
	apiv1.UnimplementedUserServiceServer
	logger *logx.Logger
}

func NewServer() *Server {
	return &Server{
		logger: logx.NewPkg("user"),
	}
}
