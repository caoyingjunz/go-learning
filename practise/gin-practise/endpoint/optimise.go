package endpoint

import (
	"github.com/gin-gonic/gin"

	"go-learning/practise/gin-practise/optimise"
)

type PreInterface interface {
	optimise.OptimiseGetter
}

// practise is used to interact with features provided by the resource
type practise struct{}

func (p *practise) Optimise(c *gin.Context) optimise.OptimiseInterface {
	return optimise.NewOptimise(c)
}

// New creates a new practise
func New() PreInterface {
	return &practise{}
}

var Practise PreInterface

func Register(p PreInterface) {
	Practise = p
}
