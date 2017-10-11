package fakeaudit

import (
	"path/filepath"
)

var replicationControllerPath = filepath.Join(absPath, "fakeaudit", "test", "replicationControllers")

func CreateFakeReplicationControllerSC(namespace string) {
	yamls := []string{"fakeReplicationControllerSC1.yml", "fakeReplicationControllerSC2.yml", "fakeReplicationControllerSC3.yml", "fakeReplicationControllerSC4.yml", "fakeReplicationControllerSC5.yml"}
	createHelper(namespace, replicationControllerPath, yamls)
}

func CreateFakeReplicationControllerRunAsNonRoot(namespace string) {
	yamls := []string{"fakeReplicationControllerRANR1.yml", "fakeReplicationControllerRANR2.yml", "fakeReplicationControllerRANR2.yml", "fakeReplicationControllerRANR3.yml", "fakeReplicationControllerRANR4.yml"}
	createHelper(namespace, replicationControllerPath, yamls)
}

func CreateFakeReplicationControllerPrivileged(namespace string) {
	yamls := []string{"fakeReplicationControllerPrivileged1.yml", "fakeReplicationControllerPrivileged2.yml", "fakeReplicationControllerPrivileged3.yml", "fakeReplicationControllerPrivileged4.yml"}
	createHelper(namespace, replicationControllerPath, yamls)
}

func CreateFakeReplicationControllerReadOnlyRootFilesystem(namespace string) {
	yamls := []string{"fakeReplicationControllerRORF1.yml", "fakeReplicationControllerRORF2.yml", "fakeReplicationControllerRORF3.yml", "fakeReplicationControllerRORF4.yml"}
	createHelper(namespace, replicationControllerPath, yamls)
}

func CreateFakeReplicationControllerAutomountServiceAccountToken(namespace string) {
	yamls := []string{"fakeReplicationControllerASAT1.yml", "fakeReplicationControllerASAT2.yml", "fakeReplicationControllerASAT3.yml"}
	createHelper(namespace, replicationControllerPath, yamls)
}

func CreateFakeReplicationControllerImg(namespace string) {
	yamls := []string{"fakeReplicationControllerImg1.yml", "fakeReplicationControllerImg2.yml"}
	createHelper(namespace, replicationControllerPath, yamls)
}
