package fakeaudit

import (
	"path/filepath"
)

var deploymetsPath = filepath.Join(absPath, "fakeaudit", "test", "deployments")

func CreateFakeDeploymentSC(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentSC1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentSC2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentSC3.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentSC4.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentSC5.yml")))
}

func CreateFakeDeploymentRunAsNonRoot(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRANR1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRANR2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRANR3.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRANR4.yml")))
}

func CreateFakeDeploymentReadOnlyRootFilesystem(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRORF1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRORF2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRORF3.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentRORF4.yml")))
}

func CreateFakeDeploymentAutomountServiceAccountToken(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentASAT1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentASAT2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentASAT3.yml")))
}

func CreateFakeDeploymentImg(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentImg1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymetsPath, "fakeDeploymentImg2.yml")))
}
