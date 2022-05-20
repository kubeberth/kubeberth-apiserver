package disks

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/gin-gonic/gin"

	"github.com/kubeberth/berth-operator/api/v1alpha1"

	"github.com/kubeberth/berth-apiserver/pkg/berth"
)

type JsonDiskSourceArchiveRequest struct {
	Name string `json:"name"`
}

type JsonDiskSourceDiskRequest struct {
	Name string `json:"name"`
}

type JsonDiskSourceRequest struct {
	Archive JsonDiskSourceArchiveRequest `json:"archive"`
	Disk JsonDiskSourceDiskRequest `json:"disk"`
}

type JsonDiskRequest struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Source JsonDiskSourceRequest `json:"source"`
	//Source JsonDiskSourceRequest `json:"source" binding:"dive"`
}

func GetAllDisks(ctx *gin.Context) {
	namespace := "kubeberth"
	disks, err := berth.Clientset.Disks().Disks(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(disks.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	var ret []v1alpha1.Disk
	for _, disk := range disks.Items {
		ret = append(ret, disk)
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetDisk(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	ret, err := berth.Clientset.Disks().Disks(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

func CreateDisk(ctx *gin.Context) {
	var j JsonDiskRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid koko",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	size := j.Size
	archiveName := j.Source.Archive.Name

	disk := &v1alpha1.Disk{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.DiskSpec{
			Size: size,
			Source: &v1alpha1.DiskSource{
				Archive: &v1alpha1.DiskSourceArchive{
					Namespace: namespace,
					Name: archiveName,
				},
			},
		},
	}

	ret, err := berth.Clientset.Disks().Disks(namespace).Create(context.TODO(), disk, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, ret)
}

func UpdateDisk(ctx *gin.Context) {
	var j JsonDiskRequest
	if err := ctx.ShouldBindJSON(&j); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := j.Name
	namespace := "kubeberth"
	size := j.Size
	archiveName := j.Source.Archive.Name
	disk, err := berth.Clientset.Disks().Disks(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	spec := v1alpha1.DiskSpec{
				Size: size,
				Source: &v1alpha1.DiskSource{
					Archive: &v1alpha1.DiskSourceArchive{
						Namespace: namespace,
						Name: archiveName,
					},
				},
			}
	disk.Spec = spec

	ret, err := berth.Clientset.Disks().Disks(namespace).Update(context.TODO(), disk, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, ret)
}

func DeleteDisk(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := berth.Clientset.Disks().Disks(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
