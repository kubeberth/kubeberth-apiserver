package cloudinits

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/gin-gonic/gin"

	"github.com/kubeberth/berth-operator/api/v1alpha1"

	"github.com/kubeberth/berth-apiserver/pkg/berth"
)

type JsonCloudInitRequest struct {
	Name string `json:"name"`
	UserData string `json:"userData"`
	NetworkData string `json:"networkData"`
}

func GetAllCloudInits(ctx *gin.Context) {
	namespace := "kubeberth"
	cloudinits, err := berth.Clientset.CloudInits().CloudInits(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(cloudinits.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	var ret []v1alpha1.CloudInit
	for _, cloudinit := range cloudinits.Items {
		ret = append(ret, cloudinit)
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetCloudInit(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	ret, err := berth.Clientset.CloudInits().CloudInits(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func CreateCloudInit(ctx *gin.Context) {
	var j JsonCloudInitRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	userData := j.UserData
	networkData := j.NetworkData

	cloudinit := &v1alpha1.CloudInit{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.CloudInitSpec{
			UserData: userData,
			NetworkData: networkData,
		},
	}

	ret, err := berth.Clientset.CloudInits().CloudInits(namespace).Create(context.TODO(), cloudinit, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func UpdateCloudInit(ctx *gin.Context) {
	var j JsonCloudInitRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	userData := j.UserData
	networkData := j.NetworkData
	cloudinit, err := berth.Clientset.CloudInits().CloudInits(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	spec := v1alpha1.CloudInitSpec{
				UserData: userData,
				NetworkData: networkData,
			}
	cloudinit.Spec = spec

	ret, err := berth.Clientset.CloudInits().CloudInits(namespace).Update(context.TODO(), cloudinit, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func DeleteCloudInit(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := berth.Clientset.CloudInits().CloudInits(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
