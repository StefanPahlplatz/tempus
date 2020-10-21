package core

import (
	"github.com/StefanPahlplatz/tempus/auth"
	"github.com/sirupsen/logrus"
)

type service struct {
	logger *logrus.Entry
}

func NewService() auth.Service {
	return &service{logger: nil}
}

func (s *service) Login(username string, password string) error {
	return nil
}
