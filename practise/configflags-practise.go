package main

import "fmt"

// ConfigFlags composes the set of values necessary
// for obtaining a REST client config
type ConfigFlags struct {
	// config flags
	ClusterName *string
	Namespace   *string

	discoveryBurst int
	discoveryQPS   float32
}

func stringptr(val string) *string {
	return &val
}

func NewConfigFlags() *ConfigFlags {
	return &ConfigFlags{
		ClusterName: stringptr(""),
		Namespace:   stringptr(""),
	}
}

func (f *ConfigFlags) WithDiscoveryBurst(discoveryBurst int) *ConfigFlags {
	f.discoveryBurst = discoveryBurst
	return f
}

func (f *ConfigFlags) WithDiscoveryQPS(discoveryQPS float32) *ConfigFlags {
	f.discoveryQPS = discoveryQPS
	return f
}

func main() {
	defaultConfigFlags := NewConfigFlags().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	fmt.Println(fmt.Sprintf("%+v", defaultConfigFlags))
}
