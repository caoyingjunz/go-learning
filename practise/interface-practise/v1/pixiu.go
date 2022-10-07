package v1

import "fmt"

// PixiuGetter has a method to return a PixiuInterface.
// A group's client should implement this interface.
type PixiuGetter interface {
	Pixiu(namespace string) PixiuInterface
}

// PixiuInterface has methods to work with Pixiu resources.
type PixiuInterface interface {
	Create(name string) error
}

// pixiu implements PixiuInterface
type pixiu struct {
	ns  string
	svc string
}

func (p *pixiu) Create(name string) error {
	fmt.Println("pixiu create", p.ns, p.svc, name)
	return nil
}

func NewPixiu(svc, namespace string) PixiuInterface {
	return &pixiu{
		ns:  namespace,
		svc: svc,
	}
}
