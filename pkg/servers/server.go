package servers

import (
	"context"
	"net/http"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gin-gonic/gin"

	"github.com/kubeberth/kubeberth-apiserver/pkg/berth"
	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type ResponseServer struct {
	Name       string                   `json:"name"`
	State      string                   `json:"state"`
	Running    bool                     `json:"running"`
	CPU        *resource.Quantity       `json:"cpu"`
	Memory     *resource.Quantity       `json:"memory"`
	MACAddress string                   `json:"mac_address"`
	IP         string                   `json:"ip"`
	Hostname   string                   `json:"hostname"`
	Hosting    string                   `json:"hosting"`
	Disks      []berth.AttachedDisk     `json:"disks"`
	ISOImage   *berth.AttachedISOImage  `json:"isoimage"`
	CloudInit  *berth.AttachedCloudInit `json:"cloudinit"`
}

type RequestServer struct {
	Name       string                   `json:"name"         binding:"required"`
	Running    bool                     `json:"running"`
	CPU        *resource.Quantity       `json:"cpu"          binding:"required"`
	Memory     *resource.Quantity       `json:"memory"       binding:"required"`
	MACAddress string                   `json:"mac_address"`
	Hostname   string                   `json:"hostname"     binding:"required"`
	Hosting    string                   `json:"hosting"`
	IP         string                   `json:"ip"`
	Disks      []berth.AttachedDisk     `json:"disks"`
	ISOImage   *berth.AttachedISOImage  `json:"isoimage"`
	CloudInit  *berth.AttachedCloudInit `json:"cloudinit"`
}

func convertServer2ResponseServer(server v1alpha1.Server) *ResponseServer {
	ret := &ResponseServer{
		Name:       server.ObjectMeta.Name,
		State:      server.Status.State,
		Running:    *server.Spec.Running,
		CPU:        server.Spec.CPU,
		Memory:     server.Spec.Memory,
		MACAddress: server.Spec.MACAddress,
		IP:         server.Status.IP,
		Hostname:   server.Spec.Hostname,
		Hosting:    server.Spec.Hosting,
		Disks:      []berth.AttachedDisk{},
		ISOImage:   &berth.AttachedISOImage{},
		CloudInit:  &berth.AttachedCloudInit{},
	}

	if ret.Hosting == "" {
		ret.Hosting = server.Status.Hosting
	}

	if server.Spec.Disks != nil {
		ret.Disks = server.Spec.Disks
	}

	if server.Spec.ISOImage != nil {
		ret.ISOImage.Name = server.Spec.ISOImage.Name
	}

	if server.Spec.CloudInit != nil {
		ret.CloudInit.Name = server.Spec.CloudInit.Name
	}

	return ret
}

func GetAllServers(ctx *gin.Context) {
	namespace := "kubeberth"
	servers, err := client.Clientset.Servers().Servers(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	var ret []*ResponseServer
	for _, server := range servers.Items {
		ret = append(ret, convertServer2ResponseServer(server))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetServer(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	server, err := client.Clientset.Servers().Servers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertServer2ResponseServer(*server))
}

func CreateServer(ctx *gin.Context) {
	var s RequestServer
	if err := ctx.ShouldBindJSON(&s); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name       := s.Name
	namespace  := "kubeberth"
	running    := s.Running
	cpu        := s.CPU
	memory     := s.Memory
	macAddress := s.MACAddress
	hostname   := s.Hostname
	hosting    := s.Hosting
	disks      := s.Disks
	isoimage   := s.ISOImage
	cloudinit  := s.CloudInit

	server := &v1alpha1.Server{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ServerSpec{
			Running:    &running,
			CPU:        cpu,
			Memory:     memory,
			MACAddress: macAddress,
			Hostname:   hostname,
			Hosting:    hosting,
			Disks:      disks,
			ISOImage:   isoimage,
			CloudInit:  cloudinit,
		},
	}

	ret, err := client.Clientset.Servers().Servers(namespace).Create(context.TODO(), server, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, convertServer2ResponseServer(*ret))
}

func UpdateServer(ctx *gin.Context) {
	var s RequestServer
	if err := ctx.ShouldBindJSON(&s); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name       := s.Name
	namespace  := "kubeberth"
	running    := s.Running
	cpu        := s.CPU
	memory     := s.Memory
	macAddress := s.MACAddress
	hostname   := s.Hostname
	hosting    := s.Hosting
	disks      := s.Disks
	isoimage   := s.ISOImage
	cloudinit  := s.CloudInit

	server, err := client.Clientset.Servers().Servers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	spec := v1alpha1.ServerSpec{
		Running:    &running,
		CPU:        cpu,
		Memory:     memory,
		MACAddress: macAddress,
		Hostname:   hostname,
		Hosting:    hosting,
		Disks:      disks,
		ISOImage:   isoimage,
		CloudInit:  cloudinit,
	}

	server.Spec = spec

	ret, err := client.Clientset.Servers().Servers(namespace).Update(context.TODO(), server, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, convertServer2ResponseServer(*ret))
}

func DeleteServer(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := client.Clientset.Servers().Servers(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
