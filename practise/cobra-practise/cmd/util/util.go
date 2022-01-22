package util

import (
	"sync"
)

type Factory interface {
	Client() string
}

type factoryImpl struct {
	client string

	// Caches OpenAPI document and parsed resources
	parser sync.Once
	getter sync.Once
}

func (f *factoryImpl) Client() string {
	return f.client
}

func NewFactory(kubeconfig string) Factory {
	c := kubeconfig
	f := &factoryImpl{
		client: c,
	}

	return f
}
