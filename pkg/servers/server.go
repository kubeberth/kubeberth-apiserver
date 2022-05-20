package servers

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/gin-gonic/gin"

	"github.com/kubeberth/berth-operator/api/v1alpha1"

	"github.com/kubeberth/berth-apiserver/pkg/berth"
)

type JsonDiskSource struct {
	Name string `json:"name"`
}

type JsonCloudInitSource struct {
	Name string `json:"name"`
}

type JsonServerRequest struct {
	Name string `json:"name"`
	Running bool `json:"running"`
	CPU string `json:"cpu"`
	Memory string `json:"memory"`
	MACAddress string `json:"macAddress"`
	HostName string `json:"hostname"`
	Disk JsonDiskSource `json:"disk"`
	CloudInit JsonCloudInitSource `json:"cloudinit"`
}

func GetAllServers(ctx *gin.Context) {
	namespace := "kubeberth"
	servers, err := berth.Clientset.Servers().Servers(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(servers.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	var ret []v1alpha1.Server
	for _, server := range servers.Items {
		ret = append(ret, server)
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetServer(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	ret, err := berth.Clientset.Servers().Servers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func CreateServer(ctx *gin.Context) {
	var j JsonServerRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid koko",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	running := j.Running
	cpu := resource.MustParse(j.CPU)
	memory := resource.MustParse(j.Memory)
	macAddress := j.MACAddress
	hostname := j.HostName
	disk := j.Disk.Name
	cloudinit := j.CloudInit.Name
	server := &v1alpha1.Server{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ServerSpec{
			Running: &running,
			CPU: &cpu,
			Memory: &memory,
			MACAddress: macAddress,
			HostName: hostname,
			Disk: &v1alpha1.DiskSourceDisk{
				Namespace: namespace,
				Name: disk,
			},
			CloudInit: &v1alpha1.CloudInitSource{
				Namespace: namespace,
				Name: cloudinit,
			},
		},
	}

	ret, err := berth.Clientset.Servers().Servers(namespace).Create(context.TODO(), server, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, ret)
}

func UpdateServer(ctx *gin.Context) {
	var j JsonServerRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid koko",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	running := j.Running
	cpu := resource.MustParse(j.CPU)
	memory := resource.MustParse(j.Memory)
	macAddress := j.MACAddress
	hostname := j.HostName
	disk := j.Disk.Name
	cloudinit := j.CloudInit.Name
	server, err := berth.Clientset.Servers().Servers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	spec := v1alpha1.ServerSpec{
			Running: &running,
			CPU: &cpu,
			Memory: &memory,
			MACAddress: macAddress,
			HostName: hostname,
			Disk: &v1alpha1.DiskSourceDisk{
				Namespace: namespace,
				Name: disk,
			},
			CloudInit: &v1alpha1.CloudInitSource{
				Namespace: namespace,
				Name: cloudinit,
			},
		}
	server.Spec = spec

	ret, err := berth.Clientset.Servers().Servers(namespace).Update(context.TODO(), server, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, ret)
}

func DeleteServer(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := berth.Clientset.Servers().Servers(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
