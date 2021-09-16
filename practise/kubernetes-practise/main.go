package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// more infomation, refer to https://github.com/kubernetes/client-go
// about the version, ecommend using the v0.x.y tags for Kubernetes releases >= v1.17.0
// and kubernetes-1.x.y tags for Kubernetes releases < v1.17.0

var clientset *kubernetes.Clientset

const (
	INGSNAME  string = "test-ingress"
	NAMESPACE string = "default"
)

func init() {
	config, err := clientcmd.BuildConfigFromFlags("", "./admins/test.conf")
	if err != nil {
		panic(err)
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

type Ingress struct {
}

type Service struct {
}

type Deployment struct {
}

func NewDploymet() *Deployment {
	return &Deployment{}
}

func (i Ingress) Get(namespace string, name string) (ing *v1beta1.Ingress, err error) {
	ing, err = clientset.ExtensionsV1beta1().
		Ingresses(namespace).
		Get(name, metav1.GetOptions{})
	return
}

func (i Ingress) Create(namespace string, name string) (ing *v1beta1.Ingress, err error) {
	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"t": "sssss",
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: "ttt.com",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Backend: v1beta1.IngressBackend{
										ServiceName: "test",
										ServicePort: intstr.FromInt(80),
									},
									Path: "/test",
								},
							},
						},
					},
				},
			},
		},
	}

	ing, err = clientset.ExtensionsV1beta1().
		Ingresses(namespace).
		Create(ingress)
	return
}

//func (i Ingress) Delete(namespace string, name string) (err error) {
//	err = clientset.ExtensionsV1beta1().
//		Ingresses(namespace).
//		Delete(name, metav1.DeleteOptions{})
//	return
//}

func (s Service) Get(namespace string, name string) (service *v1.Service, err error) {
	service, err = clientset.CoreV1().Services(namespace).
		Get(name, metav1.GetOptions{})
	return
}

func (s Service) Create(namespace string, name string) (service *v1.Service, err error) {

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"test1": "lables",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "test-port",
					Protocol:   v1.Protocol("TCP"),
					Port:       int32(80),
					TargetPort: intstr.FromInt(81),
				},
			},
		},
	}

	service, err = clientset.CoreV1().Services(namespace).
		Create(svc)
	return
}

func (s Service) Update(namespace string, name string) (service *v1.Service, err error) {
	svc, getErr := s.Get(namespace, name)
	if getErr != nil {
		return nil, getErr
	}
	svc.Spec.Ports[0].TargetPort = intstr.FromInt(98)
	service, err = clientset.CoreV1().Services(namespace).
		Update(svc)
	return
}

func (d Deployment) Get(namespace string, name string) (deploy *appsv1.Deployment, err error) {
	deploy, err = clientset.AppsV1().Deployments(namespace).
		Get(name, metav1.GetOptions{})
	return
}

func main() {

	//ingress := &Ingress{}

	service := &Service{}

	//ingGet, getErr := ingress.Get(NAMESPACE, INGSNAME)
	//if getErr != nil {
	//	panic(getErr)
	//}
	//ing, err := ingress.Create(NAMESPACE, INGSNAME)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Print(ingGet)

	//svcGet, err := service.Get(NAMESPACE, "test-ingress")
	//if err != nil {
	//	panic(err)
	//}

	//svc, err := service.Create(NAMESPACE, "test-servce")
	//if err != nil {
	//	panic(err)
	//}

	svc, err := service.Update(NAMESPACE, "test-servce")
	if err != nil {
		panic(err)
	}
	fmt.Println(svc)

	deploy := NewDploymet()
	dep, err := deploy.Get(NAMESPACE, "nginx")
	if err != nil {
		panic(err)
	}
	fmt.Println(dep)
}
