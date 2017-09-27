package fakeaudit

import (
	"path/filepath"
)

var deploymentsPath = filepath.Join(absPath, "fakeaudit", "test", "deployments")

func CreateFakeDeploymentSC(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentSC1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentSC2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentSC3.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentSC4.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentSC5.yml")))
}

func CreateFakeDeploymentRunAsNonRoot(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRANR1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRANR2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRANR3.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRANR4.yml")))
}

func CreateFakeDeploymentPrivileged(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentPrivileged1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentPrivileged2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentPrivileged3.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentPrivileged4.yml")))
}

func CreateFakeDeploymentReadOnlyRootFilesystem(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRORF1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRORF2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRORF3.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentRORF4.yml")))
}

func CreateFakeDeploymentAutomountServiceAccountToken(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentASAT1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentASAT2.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentASAT3.yml")))
}

func CreateFakeDeploymentImg(namespace string) {
	fakeDeploymentClient := getFakeDeploymentClient(namespace)
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentImg1.yml")))
	fakeDeploymentClient.Create(getDeployment(filepath.Join(deploymentsPath, "fakeDeploymentImg2.yml")))
}
