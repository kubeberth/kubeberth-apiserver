package archives

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/gin-gonic/gin"

	//"github.com/kubeberth/berth-operator/pkg/clientset/versioned/typed/"
	"github.com/kubeberth/berth-operator/api/v1alpha1"

	"github.com/kubeberth/berth-apiserver/pkg/berth"
)

type JsonArchiveRequest struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

func GetAllArchives(ctx *gin.Context) {
	namespace := "kubeberth"
	archives, err := berth.Clientset.Archives().Archives(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(archives.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	var ret []v1alpha1.Archive
	for _, archive := range archives.Items {
		ret = append(ret, archive)
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetArchive(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	ret, err := berth.Clientset.Archives().Archives(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func CreateArchive(ctx *gin.Context) {
	var j JsonArchiveRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	url := j.URL

	archive := &v1alpha1.Archive{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ArchiveSpec{
			URL: url,
		},
	}

	ret, err := berth.Clientset.Archives().Archives(namespace).Create(context.TODO(), archive, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func UpdateArchive(ctx *gin.Context) {
	var j JsonArchiveRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	url := j.URL
	archive, err := berth.Clientset.Archives().Archives(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	spec := v1alpha1.ArchiveSpec{
				URL: url,
			}
	archive.Spec = spec

	ret, err := berth.Clientset.Archives().Archives(namespace).Update(context.TODO(), archive, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func DeleteArchive(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := berth.Clientset.Archives().Archives(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
