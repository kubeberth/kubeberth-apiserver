package main

import (
	"k8s.io/klog/v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gin-gonic/gin"

	clientset "github.com/kubeberth/berth-operator/pkg/clientset/versioned"
	"github.com/kubeberth/berth-apiserver/pkg/berth"
	"github.com/kubeberth/berth-apiserver/pkg/archives"
	"github.com/kubeberth/berth-apiserver/pkg/healthz"
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

	berth.Clientset, err = clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalf("clientset.NewForConfig: %s", err.Error())
		return
	}

	g := gin.Default()
	r := g.Group("/api/v1alpha1")

	r.GET("/archives",          archives.GetAllArchives)
	r.GET("/archives/",         archives.GetAllArchives)
	r.GET("/archives/:name",    archives.GetArchive)
	r.POST("/archives",         archives.CreateArchive)
	r.POST("/archives/",        archives.CreateArchive)
	r.PUT("/archives/:name",    archives.UpdateArchive)
	r.DELETE("/archives/:name", archives.DeleteArchive)
	r.GET("/healthz",           healthz.Healthz)
	r.GET("/healthz/",          healthz.Healthz)

	klog.Info("Start")

	if err := g.Run(":2022"); err != nil {
		klog.Fatalf("start: %s", err.Error())
	}
}
