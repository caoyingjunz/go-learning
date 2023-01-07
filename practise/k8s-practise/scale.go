package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	cacheddiscovery "k8s.io/client-go/discovery/cached"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"

	"go-learning/practise/k8s-practise/app"
)

func main() {
	config, err := app.BuildClientConfig("")
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	cachedClient := cacheddiscovery.NewMemCacheClient(clientSet.Discovery())

	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(clientSet.Discovery())
	scaleClient, err := scale.NewForConfig(config, restmapper.NewDeferredDiscoveryRESTMapper(cachedClient), dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)

	gr := schema.GroupResource{
		Group:    "apps",
		Resource: "deployments",
	}

	// 1. 获取 scale, update scale
	sc, err := scaleClient.Scales("default").Get(context.TODO(), gr, "nginx", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	newSC := sc.DeepCopy()
	newSC.Spec.Replicas = newSC.Spec.Replicas + 1
	_, err = scaleClient.Scales("default").Update(context.TODO(), gr, newSC, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	//[
	//{ "op": "test", "path": "/a/b/c", "value": "foo" },
	//{ "op": "remove", "path": "/a/b/c" },
	//{ "op": "add", "path": "/a/b/c", "value": [ "foo", "bar" ] },
	//{ "op": "replace", "path": "/a/b/c", "value": 42 },
	//{ "op": "move", "from": "/a/b/c", "path": "/a/b/d" },
	//{ "op": "copy", "from": "/a/b/d", "path": "/a/b/e" }
	//]

	// 2. scale patch
	payloadTemplate := `[{ "op": "%s", "path": "/spec/replicas", "value": %s }]`
	patchPayload := fmt.Sprintf(payloadTemplate, "replace", "1")
	if _, err = scaleClient.Scales("default").Patch(context.TODO(), gr.WithVersion("v1"), "nginx", types.JSONPatchType, []byte(patchPayload), metav1.PatchOptions{}); err != nil {
		panic(err)
	}
}
