package cloudinits

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type CloudInit struct {
	Name        string `json:"name"          binding:"required"`
	UserData    string `json:"user_data"`
	NetworkData string `json:"network_data"`
}

func convertCloudInit2CloudInit(cloudinit v1alpha1.CloudInit) *CloudInit {
	ret := &CloudInit{
		Name: cloudinit.ObjectMeta.Name,
	}

	if cloudinit.Spec.UserData != "" {
		ret.UserData = cloudinit.Spec.UserData
	}

	if cloudinit.Spec.NetworkData != "" {
		ret.NetworkData = cloudinit.Spec.NetworkData
	}

	return ret
}

func GetAllCloudInits(ctx *gin.Context) {
	namespace := "kubeberth"
	cloudinits, err := client.Clientset.CloudInits().CloudInits(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(cloudinits.Items) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	var ret []*CloudInit
	for _, cloudinit := range cloudinits.Items {
		ret = append(ret, convertCloudInit2CloudInit(cloudinit))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetCloudInit(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	cloudinit, err := client.Clientset.CloudInits().CloudInits(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertCloudInit2CloudInit(*cloudinit))
}

func CreateCloudInit(ctx *gin.Context) {
	var c CloudInit
	if err := ctx.ShouldBindJSON(&c); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := c.Name
	namespace := "kubeberth"
	userData := c.UserData
	networkData := c.NetworkData

	cloudinit := &v1alpha1.CloudInit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.CloudInitSpec{
			UserData:    userData,
			NetworkData: networkData,
		},
	}

	ret, err := client.Clientset.CloudInits().CloudInits(namespace).Create(context.TODO(), cloudinit, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertCloudInit2CloudInit(*ret))
}

func UpdateCloudInit(ctx *gin.Context) {
	var c CloudInit
	if err := ctx.ShouldBindJSON(&c); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := c.Name
	namespace := "kubeberth"
	userData := c.UserData
	networkData := c.NetworkData
	cloudinit, err := client.Clientset.CloudInits().CloudInits(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	spec := v1alpha1.CloudInitSpec{
		UserData:    userData,
		NetworkData: networkData,
	}

	cloudinit.Spec = spec

	ret, err := client.Clientset.CloudInits().CloudInits(namespace).Update(context.TODO(), cloudinit, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertCloudInit2CloudInit(*ret))
}

func DeleteCloudInit(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := client.Clientset.CloudInits().CloudInits(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
