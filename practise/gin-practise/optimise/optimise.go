package optimise

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// OptimiseGetter has a method to return a OptimiseInterface.
// A group's client should implement this interface.
type OptimiseGetter interface {
	Optimise(c *gin.Context) OptimiseInterface
}

// OptimiseInterface has methods to work with Optimise resources.
type OptimiseInterface interface {
	Create(ns string) error
	Get() (string, error)

	OptimiseExpansion
}

// optimise implements OptimiseInterface
type optimise struct {
	c *gin.Context

	// 根据业务设置字段
	ns string
}

func (o *optimise) setNamespace(ns string) error {
	n := o.c.Query("namespace")
	if len(n) != 0 {
		return fmt.Errorf("empty")
	}
	o.ns = n
	return nil
}

// 处理post请求的参数
func (o *optimise) shouldBindJSON() error {

	return nil
}

type Svars struct{}

func (o *optimise) Create(ns string) error {
	var s Svars
	if err := o.c.ShouldBindJSON(&s); err != nil {
		return err
	}

	if err := o.setNamespace(ns); err != nil {
		return err
	}

	fmt.Print("optimise create", o.ns, s)
	return nil
}

func (o *optimise) Get() (string, error) {
	if err := o.setNamespace("s"); err != nil {
		return "", err
	}

	return o.ns, nil
}

func NewOptimise(c *gin.Context) *optimise {
	return &optimise{
		c: c,
	}
}

type OptimiseExpansion interface{}
