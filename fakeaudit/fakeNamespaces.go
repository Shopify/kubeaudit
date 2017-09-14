package fakeaudit

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateFakeNamespace(namespace string) (*apiv1.Namespace, error) {
	namespaceClient := getFakeNamespaceClient()
	return namespaceClient.Create(&apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
}
