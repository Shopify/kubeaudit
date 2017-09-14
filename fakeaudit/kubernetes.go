package fakeaudit

import (
	log "github.com/sirupsen/logrus"
	v1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake "k8s.io/client-go/kubernetes/fake"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

var fakeKubeClient *fake.Clientset

func init() {
	fakeKubeClient = fake.NewSimpleClientset()
}

func getFakeKubeClient() *fake.Clientset {
	return fakeKubeClient
}

func getFakeCoreClient() corev1.CoreV1Interface {
	return getFakeKubeClient().Core()
}

func getFakeAppsClient() appsv1beta1.AppsV1beta1Interface {
	return getFakeKubeClient().AppsV1beta1()
}

func getFakeExtensionClient() extensionsv1beta1.ExtensionsV1beta1Interface {
	return getFakeKubeClient().Extensions()
}

func getFakeNamespaceClient() corev1.NamespaceInterface {
	fakeCoreClient := getFakeCoreClient()
	return fakeCoreClient.Namespaces()
}

func getFakeDeploymentClient(namespace string) appsv1beta1.DeploymentInterface {
	fakeAppClient := getFakeAppsClient()
	return fakeAppClient.Deployments(namespace)
}

func getFakeStatefulSetClient(namespace string) appsv1beta1.StatefulSetInterface {
	fakeAppClient := getFakeAppsClient()
	return fakeAppClient.StatefulSets(namespace)
}

func getFakeDaemonSetClient(namespace string) extensionsv1beta1.DaemonSetInterface {
	fakeExtensionClient := getFakeExtensionClient()
	return fakeExtensionClient.DaemonSets(namespace)
}

func getFakePodClient(namespace string) corev1.PodInterface {
	fakeCoreClient := getFakeCoreClient()
	return fakeCoreClient.Pods(namespace)
}

func getFakeReplicationControllerClient(namespace string) corev1.ReplicationControllerInterface {
	fakeCoreClient := getFakeCoreClient()
	return fakeCoreClient.ReplicationControllers(namespace)
}

func GetNamespaces() *apiv1.NamespaceList {
	namespaceClient := getFakeNamespaceClient()
	namespaces, err := namespaceClient.List(metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return namespaces
}

func GetDeployments(namespace string) *v1beta1.DeploymentList {
	deploymentClient := getFakeDeploymentClient(namespace)
	deployments, err := deploymentClient.List(metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return deployments
}

func GetStatefulSets(namespace string) *v1beta1.StatefulSetList {
	statefulSetClient := getFakeStatefulSetClient(namespace)
	statefulSets, err := statefulSetClient.List(metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return statefulSets
}

func GetDaemonSets(namespace string) *extv1beta1.DaemonSetList {
	daemonSetClient := getFakeDaemonSetClient(namespace)
	daemonSets, err := daemonSetClient.List(metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return daemonSets
}

func GetPods(namespace string) *apiv1.PodList {
	podClient := getFakePodClient(namespace)
	pods, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return pods
}

func GetReplicationControllers(namespace string) *apiv1.ReplicationControllerList {
	replicationControllerClient := getFakeReplicationControllerClient(namespace)
	replicationControllers, err := replicationControllerClient.List(metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return replicationControllers
}
