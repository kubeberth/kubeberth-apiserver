package archives

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type Archive struct {
	Name       string `json:"name"       binding:"required"`
	Repository string `json:"repository"`
}

func convertArchive2Archive(archive v1alpha1.Archive) *Archive {
	ret := &Archive{
		Name:       archive.ObjectMeta.Name,
		Repository: archive.Spec.Repository,
	}

	return ret
}

func GetAllArchives(ctx *gin.Context) {
	namespace := "kubeberth"
	archives, err := client.Clientset.Archives().Archives(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(archives.Items) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	var ret []*Archive
	for _, archive := range archives.Items {
		ret = append(ret, convertArchive2Archive(archive))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetArchive(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	archive, err := client.Clientset.Archives().Archives(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertArchive2Archive(*archive))
}

func CreateArchive(ctx *gin.Context) {
	var a Archive
	if err := ctx.ShouldBindJSON(&a); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := a.Name
	namespace := "kubeberth"
	repository := a.Repository

	archive := &v1alpha1.Archive{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ArchiveSpec{
			Repository: repository,
		},
	}

	ret, err := client.Clientset.Archives().Archives(namespace).Create(context.TODO(), archive, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertArchive2Archive(*ret))
}

func UpdateArchive(ctx *gin.Context) {
	var a Archive
	if err := ctx.ShouldBindJSON(&a); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := a.Name
	namespace := "kubeberth"
	repository := a.Repository
	archive, err := client.Clientset.Archives().Archives(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	spec := v1alpha1.ArchiveSpec{
		Repository: repository,
	}

	archive.Spec = spec

	ret, err := client.Clientset.Archives().Archives(namespace).Update(context.TODO(), archive, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertArchive2Archive(*ret))
}

func DeleteArchive(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := client.Clientset.Archives().Archives(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
