package main

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/gin-gonic/gin"

	"github.com/kubeberth/kubeberth-apiserver/pkg/archives"
	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-apiserver/pkg/cloudinits"
	"github.com/kubeberth/kubeberth-apiserver/pkg/disks"
	"github.com/kubeberth/kubeberth-apiserver/pkg/healthz"
	"github.com/kubeberth/kubeberth-apiserver/pkg/loadbalancers"
	"github.com/kubeberth/kubeberth-apiserver/pkg/servers"
	clientset "github.com/kubeberth/kubeberth-operator/pkg/clientset/versioned"
)

func main() {
	klog.InitFlags(nil)
	config, err := rest.InClusterConfig()

	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			klog.Fatalf("building kubeconfig: %s", err.Error())
		}
	}

	client.Clientset, err = clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalf("clientset.NewForConfig: %s", err.Error())
		return
	}

	g := gin.Default()
	r := g.Group("/api/v1alpha1")

	r.GET("/archives", archives.GetAllArchives)
	r.GET("/archives/", archives.GetAllArchives)
	r.GET("/archives/:name", archives.GetArchive)
	r.POST("/archives", archives.CreateArchive)
	r.POST("/archives/", archives.CreateArchive)
	r.PUT("/archives/:name", archives.UpdateArchive)
	r.DELETE("/archives/:name", archives.DeleteArchive)

	r.GET("/cloudinits", cloudinits.GetAllCloudInits)
	r.GET("/cloudinits/", cloudinits.GetAllCloudInits)
	r.GET("/cloudinits/:name", cloudinits.GetCloudInit)
	r.POST("/cloudinits", cloudinits.CreateCloudInit)
	r.POST("/cloudinits/", cloudinits.CreateCloudInit)
	r.PUT("/cloudinits/:name", cloudinits.UpdateCloudInit)
	r.DELETE("/cloudinits/:name", cloudinits.DeleteCloudInit)

	r.GET("/disks", disks.GetAllDisks)
	r.GET("/disks/", disks.GetAllDisks)
	r.GET("/disks/:name", disks.GetDisk)
	r.POST("/disks", disks.CreateDisk)
	r.POST("/disks/", disks.CreateDisk)
	r.PUT("/disks/:name", disks.UpdateDisk)
	r.DELETE("/disks/:name", disks.DeleteDisk)

	r.GET("/servers", servers.GetAllServers)
	r.GET("/servers/", servers.GetAllServers)
	r.GET("/servers/:name", servers.GetServer)
	r.POST("/servers", servers.CreateServer)
	r.POST("/servers/", servers.CreateServer)
	r.PUT("/servers/:name", servers.UpdateServer)
	r.DELETE("/servers/:name", servers.DeleteServer)

	r.GET("/loadbalancers", loadbalancers.GetAllLoadBalancers)
	r.GET("/loadbalancers/", loadbalancers.GetAllLoadBalancers)
	r.GET("/loadbalancers/:name", loadbalancers.GetLoadBalancer)
	r.POST("/loadbalancers", loadbalancers.CreateLoadBalancer)
	r.POST("/loadbalancers/", loadbalancers.CreateLoadBalancer)
	r.PUT("/loadbalancers/:name", loadbalancers.UpdateLoadBalancer)
	r.DELETE("/loadbalancers/:name", loadbalancers.DeleteLoadBalancer)

	r.GET("/healthz", healthz.Healthz)
	r.GET("/healthz/", healthz.Healthz)

	klog.Info("Start")

	if err := g.Run(":2022"); err != nil {
		klog.Fatalf("start: %s", err.Error())
	}
}
