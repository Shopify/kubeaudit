package fakeaudit

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	v1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

var absPath, _ = filepath.Abs("../")

func ReadConfigFile(filename string) (decoded []runtime.Object, err error) {
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Error("File not found")
		return
	}
	buf_slice := bytes.Split(buf, []byte("---"))

	decoder := scheme.Codecs.UniversalDeserializer()

	for _, b := range buf_slice {
		obj, _, err := decoder.Decode(b, nil, nil)
		if err == nil && obj != nil {
			decoded = append(decoded, obj)
		}
	}
	return
}

func createHelper(namespace string, path string, yamls []string) {
	for _, yaml := range yamls {
		obj_slice, err := ReadConfigFile(filepath.Join(path, yaml))
		if err != nil {
			return
		}
		for _, obj := range obj_slice {
			switch resource := obj.(type) {
			case *v1beta1.Deployment:
				fakeDeploymentClient := getFakeDeploymentClient(namespace)
				fakeDeploymentClient.Create(resource)
			case *v1beta1.StatefulSet:
				fakeStatefulSetClient := getFakeStatefulSetClient(namespace)
				fakeStatefulSetClient.Create(resource)
			case *extensionsv1beta1.DaemonSet:
				fakeDaemonSetClient := getFakeDaemonSetClient(namespace)
				fakeDaemonSetClient.Create(resource)
			case *apiv1.Pod:
				fakePodClient := getFakePodClient(namespace)
				resource.Status.Phase = "Running"
				fakePodClient.Create(resource)
			case *apiv1.ReplicationController:
				fakeReplicationControllerClient := getFakeReplicationControllerClient(namespace)
				fakeReplicationControllerClient.Create(resource)
			}
		}
	}
}
