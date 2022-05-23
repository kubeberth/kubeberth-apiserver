package servers

import (
	"context"
	"net/http"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/gin-gonic/gin"

	"github.com/kubeberth/berth-operator/api/v1alpha1"

	"github.com/kubeberth/berth-apiserver/pkg/berth"
)

type DiskSource struct {
	Name string `json:"name"`
}

type CloudInitSource struct {
	Name string `json:"name"`
}

type Server struct {
	Name string `json:"name"`
	Running string `json:"running"`
	CPU string `json:"cpu"`
	Memory string `json:"memory"`
	MACAddress string `json:"macAddress"`
	HostName string `json:"hostname"`
	Disk *DiskSource `json:"disk"`
	CloudInit *CloudInitSource `json:"cloudinit"`
}

func convertServer2Server(server v1alpha1.Server) *Server {
	ret := &Server{
		Name: server.ObjectMeta.Name,
		Running: strconv.FormatBool(*server.Spec.Running),
		CPU: server.Spec.CPU.String(),
		Memory: server.Spec.Memory.String(),
		MACAddress: server.Spec.MACAddress,
		HostName: server.Spec.HostName,
		Disk: &DiskSource{
			Name: server.Spec.Disk.Name,
		},
		CloudInit: &CloudInitSource{
			Name: server.Spec.CloudInit.Name,
		},
	}

	return ret
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

	var ret []*Server
	for _, server := range servers.Items {
		ret = append(ret, convertServer2Server(server))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetServer(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	server, err := berth.Clientset.Servers().Servers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, convertServer2Server(*server))
}

func CreateServer(ctx *gin.Context) {
	var s Server
	if err := ctx.ShouldBindJSON(&s); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := s.Name
	namespace := "kubeberth"
	running, _ := strconv.ParseBool(s.Running)
	cpu := resource.MustParse(s.CPU)
	memory := resource.MustParse(s.Memory)
	macAddress := s.MACAddress
	hostname := s.HostName
	disk := s.Disk.Name
	cloudinit := s.CloudInit.Name
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

	ctx.JSON(http.StatusCreated, convertServer2Server(*ret))
}

func UpdateServer(ctx *gin.Context) {
	var s Server
	if err := ctx.ShouldBindJSON(&s); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := s.Name
	namespace := "kubeberth"
	running, _ := strconv.ParseBool(s.Running)
	cpu := resource.MustParse(s.CPU)
	memory := resource.MustParse(s.Memory)
	macAddress := s.MACAddress
	hostname := s.HostName
	disk := s.Disk.Name
	cloudinit := s.CloudInit.Name
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

	ctx.JSON(http.StatusCreated, convertServer2Server(*ret))
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
