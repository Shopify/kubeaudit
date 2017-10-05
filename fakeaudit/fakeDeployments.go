package fakeaudit

import (
	"path/filepath"
)

var deploymentsPath = filepath.Join(absPath, "fakeaudit", "test", "deployments")

func CreateFakeDeploymentSC(namespace string) {
	yamls := []string{"fakeDeploymentSC1.yml", "fakeDeploymentSC2.yml", "fakeDeploymentSC3.yml", "fakeDeploymentSC4.yml", "fakeDeploymentSC5.yml"}
	createHelper(namespace, deploymentsPath, yamls)
}

func CreateFakeDeploymentRunAsNonRoot(namespace string) {
	yamls := []string{"fakeDeploymentRANR1.yml", "fakeDeploymentRANR2.yml", "fakeDeploymentRANR3.yml", "fakeDeploymentRANR4.yml"}
	createHelper(namespace, deploymentsPath, yamls)
}

func CreateFakeDeploymentPrivileged(namespace string) {
	yamls := []string{"fakeDeploymentPrivileged1.yml", "fakeDeploymentPrivileged2.yml", "fakeDeploymentPrivileged3.yml", "fakeDeploymentPrivileged4.yml"}
	createHelper(namespace, deploymentsPath, yamls)
}

func CreateFakeDeploymentReadOnlyRootFilesystem(namespace string) {
	yamls := []string{"fakeDeploymentRORF1.yml", "fakeDeploymentRORF2.yml", "fakeDeploymentRORF3.yml", "fakeDeploymentRORF4.yml"}
	createHelper(namespace, deploymentsPath, yamls)
}

func CreateFakeDeploymentAutomountServiceAccountToken(namespace string) {
	yamls := []string{"fakeDeploymentASAT1.yml", "fakeDeploymentASAT2.yml", "fakeDeploymentASAT3.yml"}
	createHelper(namespace, deploymentsPath, yamls)
}

func CreateFakeDeploymentImg(namespace string) {
	yamls := []string{"fakeDeploymentImg1.yml", "fakeDeploymentImg2.yml"}
	createHelper(namespace, deploymentsPath, yamls)
}
