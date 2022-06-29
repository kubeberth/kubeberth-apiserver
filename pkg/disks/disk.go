package disks

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeberth/kubeberth-apiserver/pkg/berth"
	"github.com/kubeberth/kubeberth-apiserver/pkg/client"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type Disk struct {
	Name   string                `json:"name"    binding:"required"`
	Size   string                `json:"size"    binding:"required"`
	Source *berth.AttachedSource `json:"source"`
}

func convertDisk2Disk(disk v1alpha1.Disk) *Disk {
	ret := &Disk{
		Name:   disk.ObjectMeta.Name,
		Size:   disk.Spec.Size,
		Source: &berth.AttachedSource{},
	}

	if disk.Spec.Source != nil {
		if disk.Spec.Source.Archive != nil {
			ret.Source.Archive = &berth.AttachedArchive{
				Name: disk.Spec.Source.Archive.Name,
			}
		}

		if disk.Spec.Source.Disk != nil {
			ret.Source.Disk = &berth.AttachedDisk{
				Name: disk.Spec.Source.Disk.Name,
			}
		}
	}

	return ret
}

func GetAllDisks(ctx *gin.Context) {
	namespace := "kubeberth"
	disks, err := client.Clientset.Disks().Disks(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, convertDisk2Disk(*disk))
}

func CreateDisk(ctx *gin.Context) {
	var d Disk
	if err := ctx.ShouldBindJSON(&d); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := d.Name
	namespace := "kubeberth"
	size := d.Size
	var source *berth.AttachedSource

	if d.Source != nil {
		source = &berth.AttachedSource{}

		if d.Source.Archive != nil {
			source.Archive = &berth.AttachedArchive{
				Name: d.Source.Archive.Name,
			}
		}

		if d.Source.Disk != nil {
			source.Disk = &berth.AttachedDisk{
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, convertDisk2Disk(*ret))
}

func UpdateDisk(ctx *gin.Context) {
	var d Disk
	if err := ctx.ShouldBindJSON(&d); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid: " + err.Error(),
		})
		return
	}

	name := d.Name
	namespace := "kubeberth"
	size := d.Size
	archiveName := d.Source.Archive.Name
	disk, err := client.Clientset.Disks().Disks(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
		})
		return
	}

	spec := v1alpha1.DiskSpec{
		Size: size,
		Source: &berth.AttachedSource{
			Archive: &berth.AttachedArchive{
				Name: archiveName,
			},
		},
	}

	disk.Spec = spec

	ret, err := client.Clientset.Disks().Disks(namespace).Update(context.TODO(), disk, metav1.UpdateOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "update error: " + err.Error(),
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
