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

type Server struct {
	Name       string                      `json:"name"`
	Running    bool                        `json:"running"`
	CPU        string                      `json:"cpu"`
	Memory     string                      `json:"memory"`
	MACAddress string                      `json:"mac_address"`
	Hostname   string                      `json:"hostname"`
	Hosting    string                      `json:"hosting"`
	Disk       *berth.AttachedDisk      `json:"disk"`
	CloudInit  *berth.AttachedCloudInit `json:"cloudinit"`
}

func convertServer2Server(server v1alpha1.Server) *Server {
	ret := &Server{
		Name:       server.ObjectMeta.Name,
		Running:    *server.Spec.Running,
		CPU:        server.Spec.CPU.String(),
		Memory:     server.Spec.Memory.String(),
		MACAddress: server.Spec.MACAddress,
		Hostname:   server.Spec.Hostname,
		Hosting:    server.Spec.Hosting,
		Disk: &berth.AttachedDisk{
			Name: server.Spec.Disk.Name,
		},
		CloudInit: &berth.AttachedCloudInit{
			Name: server.Spec.CloudInit.Name,
		},
	}

	return ret
}

func GetAllServers(ctx *gin.Context) {
	namespace := "kubeberth"
	servers, err := client.Clientset.Servers().Servers(namespace).List(context.TODO(), metav1.ListOptions{})

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
	server, err := client.Clientset.Servers().Servers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

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
	running := s.Running
	cpu := resource.MustParse(s.CPU)
	memory := resource.MustParse(s.Memory)
	macAddress := s.MACAddress
	hostname := s.Hostname
	hosting := s.Hosting
	disk := s.Disk.Name
	cloudinit := s.CloudInit.Name

	server := &v1alpha1.Server{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ServerSpec{
			Running:    &running,
			CPU:        &cpu,
			Memory:     &memory,
			MACAddress: macAddress,
			Hostname:   hostname,
			Hosting:    hosting,
			Disk: &berth.AttachedDisk{
				Name: disk,
			},
			CloudInit: &berth.AttachedCloudInit{
				Name: cloudinit,
			},
		},
	}

	ret, err := client.Clientset.Servers().Servers(namespace).Create(context.TODO(), server, metav1.CreateOptions{})
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
	running := s.Running
	cpu := resource.MustParse(s.CPU)
	memory := resource.MustParse(s.Memory)
	macAddress := s.MACAddress
	hostname := s.Hostname
	hosting := s.Hosting
	disk := s.Disk.Name
	cloudinit := s.CloudInit.Name

	server, err := client.Clientset.Servers().Servers(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	spec := v1alpha1.ServerSpec{
		Running:    &running,
		CPU:        &cpu,
		Memory:     &memory,
		MACAddress: macAddress,
		Hostname:   hostname,
		Hosting:    hosting,
		Disk: &berth.AttachedDisk{
			Name: disk,
		},
		CloudInit: &berth.AttachedCloudInit{
			Name: cloudinit,
		},
	}

	server.Spec = spec

	ret, err := client.Clientset.Servers().Servers(namespace).Update(context.TODO(), server, metav1.UpdateOptions{})
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
	err := client.Clientset.Servers().Servers(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
