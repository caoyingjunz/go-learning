package v2

import "fmt"

// ScalerGetter has a method to return a ScalerInterface.
// A group's client should implement this interface.
type ScalerGetter interface {
	Scaler(namespace string) ScalerInterface
}

// ScalerInterface has methods to work with Scaler resources.
type ScalerInterface interface {
	Create(name string) error
}

// scaler implements ScalerInterface
type scaler struct {
	ns  string
	svc string
}

func (p *scaler) Create(name string) error {
	fmt.Println("create scaler", p.ns, p.svc, name)
	return nil
}

func NewScaler(svc, namespace string) ScalerInterface {
	return &scaler{
		ns:  namespace,
		svc: svc,
	}
}
