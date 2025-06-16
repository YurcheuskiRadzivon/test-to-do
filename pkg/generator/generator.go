package generator

import "github.com/google/uuid"

type UUVGenerator interface {
	NewFileName() string
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) NewFileName() string {
	return uuid.New().String()
}
