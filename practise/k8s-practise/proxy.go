package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// 请求验证 curl 127.0.0.1:8090/apis/apps/v1/namespaces/default/deployments/toolbox

func main() {
	route := gin.Default()

	// gin 指定代理，apis 原始请求转发到 k8s APIServer
	route.Any("/apis/*proxy", proxyHandler)

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

	l := *c.Request.URL
	l.Host = "59.111.229.69:6443"
	l.Scheme = "https"

	httpProxy := proxy.NewUpgradeAwareHandler(&l, transport, true, false, nil)
	httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(transport, transport)
	httpProxy.ServeHTTP(c.Writer, c.Request)
}
