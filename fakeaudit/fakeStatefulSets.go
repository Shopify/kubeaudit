package fakeaudit

import (
	"path/filepath"
)

var statefulSetPath = filepath.Join(absPath, "fakeaudit", "test", "statefulSets")

func CreateFakeStatefulSetSC(namespace string) {
	yamls := []string{"fakeStatefulSetSC2.yml", "fakeStatefulSetSC1.yml", "fakeStatefulSetSC3.yml", "fakeStatefulSetSC4.yml", "fakeStatefulSetSC5.yml"}
	createHelper(namespace, statefulSetPath, yamls)
}

func CreateFakeStatefulSetRunAsNonRoot(namespace string) {
	yamls := []string{"fakeStatefulSetRANR1.yml", "fakeStatefulSetRANR2.yml", "fakeStatefulSetRANR3.yml", "fakeStatefulSetRANR4.yml"}
	createHelper(namespace, statefulSetPath, yamls)
}

func CreateFakeStatefulSetPrivileged(namespace string) {
	yamls := []string{"fakeStatefulSetPrivileged1.yml", "fakeStatefulSetPrivileged2.yml", "fakeStatefulSetPrivileged3.yml", "fakeStatefulSetPrivileged4.yml"}
	createHelper(namespace, statefulSetPath, yamls)
}

func CreateFakeStatefulSetReadOnlyRootFilesystem(namespace string) {
	yamls := []string{"fakeStatefulSetRORF1.yml", "fakeStatefulSetRORF2.yml", "fakeStatefulSetRORF3.yml", "fakeStatefulSetRORF4.yml"}
	createHelper(namespace, statefulSetPath, yamls)
}

func CreateFakeStatefulSetAutomountServiceAccountToken(namespace string) {
	yamls := []string{"fakeStatefulSetASAT1.yml", "fakeStatefulSetASAT2.yml", "fakeStatefulSetASAT3.yml"}
	createHelper(namespace, statefulSetPath, yamls)
}

func CreateFakeStatefulSetImg(namespace string) {
	yamls := []string{"fakeStatefulSetImg1.yml", "fakeStatefulSetImg2.yml"}
	createHelper(namespace, statefulSetPath, yamls)
}
