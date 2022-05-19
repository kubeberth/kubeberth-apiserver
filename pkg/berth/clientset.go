package berth

import (
	clientset "github.com/kubeberth/berth-operator/pkg/clientset/versioned"
)

var (
	Clientset *clientset.Clientset
)
