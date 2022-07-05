package berth

import (
	corev1 "k8s.io/api/core/v1"
	"github.com/kubeberth/kubeberth-operator/api/v1alpha1"
)

type AttachedArchive   = v1alpha1.AttachedArchive
type AttachedCloudInit = v1alpha1.AttachedCloudInit
type AttachedDisk      = v1alpha1.AttachedDisk
type AttachedSource    = v1alpha1.AttachedSource
type Destination       = v1alpha1.Destination
type Port              = corev1.ServicePort
