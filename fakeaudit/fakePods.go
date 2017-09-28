package fakeaudit

import (
	"path/filepath"
)

var podsPath = filepath.Join(absPath, "fakeaudit", "test", "pods")

func CreateFakePodSC(namespace string) {
	fakePodClient := getFakePodClient(namespace)
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodSC1.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodSC2.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodSC3.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodSC4.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodSC5.yml")))
}

func CreateFakePodRunAsNonRoot(namespace string) {
	fakePodClient := getFakePodClient(namespace)
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRANR1.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRANR2.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRANR3.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRANR4.yml")))
}

func CreateFakePodPrivileged(namespace string) {
	fakePodClient := getFakePodClient(namespace)
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodPrivileged1.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodPrivileged2.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodPrivileged3.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodPrivileged4.yml")))
}

func CreateFakePodReadOnlyRootFilesystem(namespace string) {
	fakePodClient := getFakePodClient(namespace)
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRORF1.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRORF2.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRORF3.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodRORF4.yml")))
}

func CreateFakePodAutomountServiceAccountToken(namespace string) {
	fakePodClient := getFakePodClient(namespace)
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodASAT1.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodASAT2.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodASAT3.yml")))
}

func CreateFakePodImg(namespace string) {
	fakePodClient := getFakePodClient(namespace)
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodImg1.yml")))
	fakePodClient.Create(getPod(filepath.Join(podsPath, "fakePodImg2.yml")))
}
