package core

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// Validation for project name
	ProjectNameRegex      = `^[a-zA-Z0-9_]+$`
	ErrInvalidProjectName = errors.New("invalid project name")
)

type Name string

type Project struct {
	Name Name
}

func NewProject(n string) (*Project, error) {
	name, err := NewName(n)
	if err != nil {
		return nil, err
	}

	return &Project{Name: name}, nil
}

func NewName(n string) (Name, error) {
	name := strings.TrimSpace(strings.ToLower(n))

	parse, err := regexp.Compile(ProjectNameRegex)
	if err != nil {
		return "", err
	}

	if !parse.MatchString(name) {
		return "", ErrInvalidProjectName
	}

	if name == "" {
		return "", ErrInvalidProjectName
	}

	return Name(name), nil
}

func (p *Project) GetName() string {
	return string(p.Name)
}
