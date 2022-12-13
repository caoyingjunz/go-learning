package main

import (
	"net/url"
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
	route.Any("/*proxy", proxyHandler)

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
	target, err := parseTarget(*c.Request.URL, config.Host)
	if err != nil {
		panic(err)
	}

	httpProxy := proxy.NewUpgradeAwareHandler(target, transport, false, false, nil)
	httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(transport, transport)
	httpProxy.ServeHTTP(c.Writer, c.Request)
}

func parseTarget(target url.URL, host string) (*url.URL, error) {
	kubeURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	target.Host = kubeURL.Host
	target.Scheme = kubeURL.Scheme
	return &target, nil
}
