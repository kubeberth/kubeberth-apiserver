package disks

import (
	"context"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/gin-gonic/gin"

	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"

	"github.com/kubeberth/kubeberth-apiserver/pkg/berth"
)

type Disk struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Source *berth.AttachedSource `json:"source"`
}

func convertDisk2Disk(disk v1alpha1.Disk) *Disk {
	ret := &Disk{
		Name: disk.ObjectMeta.Name,
		Size: disk.Spec.Size,
		Source: &berth.AttachedSource{},
	}

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

	return ret
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

	var ret []*Disk
	for _, disk := range disks.Items {
		ret = append(ret, convertDisk2Disk(disk))
	}

	ctx.JSON(http.StatusOK, ret)
}

func GetDisk(ctx *gin.Context) {
	name := ctx.Param("name")
	namespace := "kubeberth"
	disk, err := berth.Clientset.Disks().Disks(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, convertDisk2Disk(*disk))
}

func CreateDisk(ctx *gin.Context) {
	var d Disk
	if err := ctx.ShouldBindJSON(&d); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "request invalid koko",
		})
		return
	}

	name := d.Name
	namespace := "kubeberth"
	size := d.Size
	archiveName := d.Source.Archive.Name

	disk := &v1alpha1.Disk{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.DiskSpec{
			Size: size,
			Source: &v1alpha1.AttachedSource{
				Archive: &v1alpha1.AttachedArchive{
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
	disk, err := berth.Clientset.Disks().Disks(namespace).Get(context.TODO(), name, metav1.GetOptions{})

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

	ctx.JSON(http.StatusCreated, convertDisk2Disk(*ret))
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
