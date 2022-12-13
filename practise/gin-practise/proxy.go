package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	route := gin.Default()

	// gin 制作代理，原始请求转发到 k8s APIServer
	route.Any("/apis/*action", proxyHandler)

	_ = route.Run(":8090")
}

func proxyHandler(c *gin.Context) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}

	transport, err := rest.TransportFor(config)
	if err != nil {
		panic(err)
	}

	s := *c.Request.URL
	s.Host = "175.102.24.135:6443"
	s.Scheme = "https"

	httpProxy := proxy.NewUpgradeAwareHandler(&s, transport, true, false, nil)
	httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(transport, transport)
	httpProxy.ServeHTTP(c.Writer, c.Request)
}
