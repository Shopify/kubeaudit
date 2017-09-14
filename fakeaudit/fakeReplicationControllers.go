package fakeaudit

import (
	"path/filepath"
)

var replicationControllerPath = filepath.Join(absPath, "fakeaudit", "test", "replicationControllers")

func CreateFakeReplicationControllerSC(namespace string) {
	fakeReplicationControllerClient := getFakeReplicationControllerClient(namespace)
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerSC1.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerSC2.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerSC3.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerSC4.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerSC5.yml")))
}

func CreateFakeReplicationControllerRunAsNonRoot(namespace string) {
	fakeReplicationControllerClient := getFakeReplicationControllerClient(namespace)
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRANR1.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRANR2.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRANR3.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRANR4.yml")))
}

func CreateFakeReplicationControllerReadOnlyRootFilesystem(namespace string) {
	fakeReplicationControllerClient := getFakeReplicationControllerClient(namespace)
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRORF1.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRORF2.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRORF3.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerRORF4.yml")))
}

func CreateFakeReplicationControllerAutomountServiceAccountToken(namespace string) {
	fakeReplicationControllerClient := getFakeReplicationControllerClient(namespace)
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerASAT1.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerASAT2.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerASAT3.yml")))
}

func CreateFakeReplicationControllerImg(namespace string) {
	fakeReplicationControllerClient := getFakeReplicationControllerClient(namespace)
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerImg1.yml")))
	fakeReplicationControllerClient.Create(getReplicationController(filepath.Join(replicationControllerPath, "fakeReplicationControllerImg2.yml")))
}
