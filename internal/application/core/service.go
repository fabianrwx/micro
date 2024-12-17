package core

import (
	"errors"
	"strings"
)

var (
	ErrInvalidServiceName = errors.New("invalid service name")
)

type ServiceName string
type ModuleName string

type Service struct {
	Name       ServiceName
	ModuleName ModuleName
}

func (s *Service) GetName() string {
	return string(s.Name)
}
func NewService(name, module string) (*Service, error) {
	sn, err := NewServiceName(name)
	if err != nil {
		return nil, err
	}

	mn, err := NewModuleName(module)
	if err != nil {
		return nil, err
	}

	return &Service{Name: sn, ModuleName: mn}, nil
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

func (s ModuleName) GetModuleName() string {
	return string(s)
}

func NewModuleName(name string) (ModuleName, error) {
	sn := strings.TrimSpace(strings.ToLower(name))

	if sn == "" {
		return "", ErrInvalidServiceName
	}

	return ModuleName(sn), nil
}
