package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var ctx = context.Background()

func GetNamespace(id int) ([]v1.Namespace, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}
	namespaceList, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return namespaceList.Items, err
}

//	func GetNamespaceSub(id int,nsName string) ([]k8s.NameSpace, error) {
//		clientSet, err := k8s.GetKubeConfig(id)
//		if err != nil {
//			fmt.Println(err)
//		}
//		namespacesub,err := clientSet.CoreV1().Namespaces().Get(ctx,nsName,metav1.GetOptions{})
//
//		if err != nil {
//			klog.Error(err)
//			return nil, err
//		}
//		return namespacesub, err
//	}
func AddNamespace(id int, ns k8s.NameSpace) (*v1.Namespace, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
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
func EditNamespace(id int, ns k8s.NameSpace) (*v1.Namespace, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}
	newNamespace, err := clientSet.CoreV1().Namespaces().Update(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ns.Name,
			Labels:      ns.Labels,
			Annotations: ns.Annotations,
		},
	}, metav1.UpdateOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return newNamespace, nil
}

func DelNamespace(id int, nsName string) (code int, err error) {
	clientSet, err := k8s.GetKubeConfig(id)
	deletePolicy := metav1.DeletePropagationForeground
	err = clientSet.CoreV1().Namespaces().Delete(ctx, nsName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return 400, err
	}

	return 200, err
}
