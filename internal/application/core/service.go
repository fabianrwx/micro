package core

import (
	"errors"
	"strings"
)

var (
	ErrInvalidServiceName = errors.New("invalid service name")
)

type ServiceName string

type Service struct {
	Name ServiceName
}

func (s *Service) GetName() string {
	return string(s.Name)
}
func NewService(name string) (*Service, error) {
	sn, err := NewServiceName(name)
	if err != nil {
		return nil, err
	}

	return &Service{Name: sn}, nil
}

func NewServiceName(name string) (ServiceName, error) {
	sn := strings.TrimSpace(strings.ToLower(name))

	if sn == "" {
		return "", ErrInvalidServiceName
	}

	return ServiceName(sn), nil
}

func (s ServiceName) GetServiceName() string {
	return string(s)
}
