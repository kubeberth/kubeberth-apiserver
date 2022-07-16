package loadbalancers

import (
	"context"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gin-gonic/gin"

	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type ResponseLoadBalancer struct {
	Name           string                 `json:"name"`
	State          string                 `json:"state"`
	IP             string                 `json:"ip"`
	Backends       []v1alpha1.Destination `json:"backends"`
	Ports          []corev1.ServicePort   `json:"ports"`
	BackendsStatus map[string]string      `json:backendsStatus"`
	Health         string                 `json:"health"`
}

type RequestLoadBalancer struct {
	Name     string                 `json:"name"     binding:"required"`
	Backends []v1alpha1.Destination `json:"backends" binding:"required"`
	Ports    []corev1.ServicePort   `json:"ports"    binding:"required"`
}

func convertLoadBalancer2ResponseLoadBalancer(loadbalancer v1alpha1.LoadBalancer) *ResponseLoadBalancer {
	ret := &ResponseLoadBalancer{
		Name:     loadbalancer.GetName(),
		State:    loadbalancer.Status.State,
		IP:       loadbalancer.Status.IP,
		Backends: loadbalancer.Status.Backends,
		Ports:    loadbalancer.Spec.Ports,
		BackendsStatus: loadbalancer.Status.BackendsStatus,
		Health:   loadbalancer.Status.Health,
	}

	return ret
}

func GetAllLoadBalancers(ctx *gin.Context) {
	namespace := "kubeberth"
	loadbalancers, err := client.Clientset.LoadBalancers().LoadBalancers(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	var ret []*ResponseLoadBalancer
	for _, loadbalancer := range loadbalancers.Items {
		ret = append(ret, convertLoadBalancer2ResponseLoadBalancer(loadbalancer))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetLoadBalancer(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	loadbalancer, err := client.Clientset.LoadBalancers().LoadBalancers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertLoadBalancer2ResponseLoadBalancer(*loadbalancer))
}

func CreateLoadBalancer(ctx *gin.Context) {
	var lb RequestLoadBalancer
	if err := ctx.ShouldBindJSON(&lb); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := lb.Name
	namespace := "kubeberth"
	backends := lb.Backends
	ports := lb.Ports

	loadbalancer := &v1alpha1.LoadBalancer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.LoadBalancerSpec{
			Backends: backends,
			Ports:    ports,
		},
	}

	ret, err := client.Clientset.LoadBalancers().LoadBalancers(namespace).Create(context.TODO(), loadbalancer, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, convertLoadBalancer2ResponseLoadBalancer(*ret))
}

func UpdateLoadBalancer(ctx *gin.Context) {
	var lb RequestLoadBalancer
	if err := ctx.ShouldBindJSON(&lb); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := lb.Name
	namespace := "kubeberth"
	backends := lb.Backends
	ports := lb.Ports

	loadbalancer, err := client.Clientset.LoadBalancers().LoadBalancers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	spec := v1alpha1.LoadBalancerSpec{
		Backends: backends,
		Ports:    ports,
	}

	loadbalancer.Spec = spec

	ret, err := client.Clientset.LoadBalancers().LoadBalancers(namespace).Update(context.TODO(), loadbalancer, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, convertLoadBalancer2ResponseLoadBalancer(*ret))
}

func DeleteLoadBalancer(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := client.Clientset.LoadBalancers().LoadBalancers(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
