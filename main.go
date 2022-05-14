package main

import (
	"k8s.io/klog/v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gin-gonic/gin"

	"github.com/kubeberth/berth-apiserver/pkg/healthz"
)

func main() {
	klog.InitFlags(nil)
	config, err := rest.InClusterConfig()

	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			klog.Fatalf("Error building kubeconfig: %s", err.Error())
		}
	}

	klog.Info(config)

	r := gin.Default()

	r.GET("/api/v1alpha1/healthz", healthz.Healthz)
	r.GET("/api/v1alpha1/healthz/", healthz.Healthz)

	klog.Info("start")

	if err := r.Run(":2022"); err != nil {
		klog.Fatalf(err.Error())
	}
}
