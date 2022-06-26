package disks

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type Disk struct {
	Name   string                   `json:"name"`
	Size   string                   `json:"size"`
	Source *v1alpha1.AttachedSource `json:"source"`
}

func convertDisk2Disk(disk v1alpha1.Disk) *Disk {
	ret := &Disk{
		Name:   disk.ObjectMeta.Name,
		Size:   disk.Spec.Size,
		Source: &v1alpha1.AttachedSource{},
	}

	if disk.Spec.Source != nil {
		if disk.Spec.Source.Archive != nil {
			ret.Source.Archive = &v1alpha1.AttachedArchive{
				Name: disk.Spec.Source.Archive.Name,
			}
		}

		if disk.Spec.Source.Disk != nil {
			ret.Source.Disk = &v1alpha1.AttachedDisk{
				Name: disk.Spec.Source.Disk.Name,
			}
		}
	}

	return ret
}

func GetAllDisks(ctx *gin.Context) {
	namespace := "kubeberth"
	disks, err := client.Clientset.Disks().Disks(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil || len(disks.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	var ret []*Disk
	for _, disk := range disks.Items {
		ret = append(ret, convertDisk2Disk(disk))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetDisk(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	disk, err := client.Clientset.Disks().Disks(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, convertDisk2Disk(*disk))
}

func CreateDisk(ctx *gin.Context) {
	var d Disk
	if err := ctx.ShouldBindJSON(&d); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := d.Name
	namespace := "kubeberth"
	size := d.Size
	var source *v1alpha1.AttachedSource

	if d.Source != nil {
		source = &v1alpha1.AttachedSource{}

		if d.Source.Archive != nil {
			source.Archive = &v1alpha1.AttachedArchive{
				Name: d.Source.Archive.Name,
			}
		}

		if d.Source.Disk != nil {
			source.Disk = &v1alpha1.AttachedDisk{
				Name: d.Source.Disk.Name,
			}
		}
	}

	disk := &v1alpha1.Disk{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.DiskSpec{
			Size:   size,
			Source: source,
		},
	}

	ret, err := client.Clientset.Disks().Disks(namespace).Create(context.TODO(), disk, metav1.CreateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, convertDisk2Disk(*ret))
}

func UpdateDisk(ctx *gin.Context) {
	var d Disk
	if err := ctx.ShouldBindJSON(&d); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid",
		})
		return
	}

	name := d.Name
	namespace := "kubeberth"
	size := d.Size
	archiveName := d.Source.Archive.Name
	disk, err := client.Clientset.Disks().Disks(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	spec := v1alpha1.DiskSpec{
		Size: size,
		Source: &v1alpha1.AttachedSource{
			Archive: &v1alpha1.AttachedArchive{
				Name: archiveName,
			},
		},
	}

	disk.Spec = spec

	ret, err := client.Clientset.Disks().Disks(namespace).Update(context.TODO(), disk, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "update error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, convertDisk2Disk(*ret))
}

func DeleteDisk(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	err := client.Clientset.Disks().Disks(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

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
