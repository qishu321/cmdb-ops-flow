package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"log"
)

var ctx = context.Background()

func GetNamespace(id int) ([]v1.Namespace, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		log.Fatal(err)
	}
	namespaceList, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return namespaceList.Items, err
}

func AddNamespace(id int, ns k8s.NameSpace) (*v1.Namespace, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		log.Fatal(err)
	}
	newNamespace, err := clientSet.CoreV1().Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ns.Name,
			Labels:      ns.Labels,
			Annotations: ns.Annotations,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return newNamespace, nil
}
