package fakeaudit

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	v1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"path/filepath"
)

var absPath, _ = filepath.Abs("../")

func readConfigFiles(filename string) runtime.Object {
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Error("File not found")
		panic(err)
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	obj, _, err := decoder.Decode(buf, nil, nil)

	if err != nil {
		log.Errorf("Could not decode the given yaml: %s\n%s", string(buf), err)
	}

	return obj
}

func getDeployment(filename string) (deployment *v1beta1.Deployment) {
	obj := readConfigFiles(filename)
	deployment = obj.(*v1beta1.Deployment)
	return
}

func getStatefulSet(filename string) (statefulSet *v1beta1.StatefulSet) {
	obj := readConfigFiles(filename)
	statefulSet = obj.(*v1beta1.StatefulSet)
	return
}

func getDaemonSet(filename string) (daemonSet *extensionsv1beta1.DaemonSet) {
	obj := readConfigFiles(filename)
	daemonSet = obj.(*extensionsv1beta1.DaemonSet)
	return
}

func getPod(filename string) (pod *apiv1.Pod) {
	obj := readConfigFiles(filename)
	pod = obj.(*apiv1.Pod)
	pod.Status.Phase = "Running"
	return
}

func getReplicationController(filename string) (rc *apiv1.ReplicationController) {
	obj := readConfigFiles(filename)
	rc = obj.(*apiv1.ReplicationController)
	return
}
