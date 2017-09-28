package fakeaudit

import (
	"path/filepath"
)

var statefulSetPath = filepath.Join(absPath, "fakeaudit", "test", "statefulSets")

func CreateFakeStatefulSetSC(namespace string) {
	fakeStatefulSetClient := getFakeStatefulSetClient(namespace)
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetSC1.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetSC2.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetSC3.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetSC4.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetSC5.yml")))
}

func CreateFakeStatefulSetRunAsNonRoot(namespace string) {
	fakeStatefulSetClient := getFakeStatefulSetClient(namespace)
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRANR1.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRANR2.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRANR3.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRANR4.yml")))
}

func CreateFakeStatefulSetPrivileged(namespace string) {
	fakeStatefulSetClient := getFakeStatefulSetClient(namespace)
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetPrivileged1.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetPrivileged2.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetPrivileged3.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetPrivileged4.yml")))
}

func CreateFakeStatefulSetReadOnlyRootFilesystem(namespace string) {
	fakeStatefulSetClient := getFakeStatefulSetClient(namespace)
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRORF1.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRORF2.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRORF3.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetRORF4.yml")))
}

func CreateFakeStatefulSetAutomountServiceAccountToken(namespace string) {
	fakeStatefulSetClient := getFakeStatefulSetClient(namespace)
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetASAT1.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetASAT2.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetASAT3.yml")))
}

func CreateFakeStatefulSetImg(namespace string) {
	fakeStatefulSetClient := getFakeStatefulSetClient(namespace)
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetImg1.yml")))
	fakeStatefulSetClient.Create(getStatefulSet(filepath.Join(statefulSetPath, "fakeStatefulSetImg2.yml")))
}
