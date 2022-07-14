package isoimages

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type ResponseISOImage struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	Size       string `json:"size"`
	Repository string `json:"repository"`
}

type RequestISOImage struct {
	Name       string `json:"name"       binding:"required"`
	Size       string `json:"size"       binding:"required"`
	Repository string `json:"repository" binding:"required"`
}

func convertISOImage2ISOImage(isoimage v1alpha1.ISOImage) *ResponseISOImage {
	ret := &ResponseISOImage{
		Name:       isoimage.ObjectMeta.Name,
		State:      isoimage.Status.State,
		Size:       isoimage.Spec.Size,
		Repository: isoimage.Spec.Repository,
	}

	return ret
}

func GetAllISOImages(ctx *gin.Context) {
	namespace := "kubeberth"
	isoimages, err := client.Clientset.ISOImages().ISOImages(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(isoimages.Items) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	var ret []*ResponseISOImage
	for _, isoimage := range isoimages.Items {
		ret = append(ret, convertISOImage2ISOImage(isoimage))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetISOImage(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	isoimage, err := client.Clientset.ISOImages().ISOImages(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertISOImage2ISOImage(*isoimage))
}

func CreateISOImage(ctx *gin.Context) {
	var iso RequestISOImage
	if err := ctx.ShouldBindJSON(&iso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := iso.Name
	namespace := "kubeberth"
	size := iso.Size
	repository := iso.Repository

	isoimage := &v1alpha1.ISOImage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ISOImageSpec{
			Size:       size,
			Repository: repository,
		},
	}

	ret, err := client.Clientset.ISOImages().ISOImages(namespace).Create(context.TODO(), isoimage, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertISOImage2ISOImage(*ret))
}

func UpdateISOImage(ctx *gin.Context) {
	var iso RequestISOImage
	if err := ctx.ShouldBindJSON(&iso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := iso.Name
	namespace := "kubeberth"
	size := iso.Size
	repository := iso.Repository
	isoimage, err := client.Clientset.ISOImages().ISOImages(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	spec := v1alpha1.ISOImageSpec{
		Size: size,
		Repository: repository,
	}

	isoimage.Spec = spec

	ret, err := client.Clientset.ISOImages().ISOImages(namespace).Update(context.TODO(), isoimage, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertISOImage2ISOImage(*ret))
}

func DeleteISOImage(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := client.Clientset.ISOImages().ISOImages(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
