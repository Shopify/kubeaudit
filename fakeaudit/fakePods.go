package fakeaudit

import (
	"path/filepath"
)

var podsPath = filepath.Join(absPath, "fakeaudit", "test", "pods")

func CreateFakePodSC(namespace string) {
	yamls := []string{"fakePodSC2.yml", "fakePodSC1.yml", "fakePodSC3.yml", "fakePodSC4.yml", "fakePodSC5.yml"}
	createHelper(namespace, podsPath, yamls)
}

func CreateFakePodRunAsNonRoot(namespace string) {
	yamls := []string{"fakePodRANR1.yml", "fakePodRANR2.yml", "fakePodRANR3.yml", "fakePodRANR4.yml"}
	createHelper(namespace, podsPath, yamls)
}

func CreateFakePodPrivileged(namespace string) {
	yamls := []string{"fakePodPrivileged1.yml", "fakePodPrivileged2.yml", "fakePodPrivileged3.yml", "fakePodPrivileged4.yml"}
	createHelper(namespace, podsPath, yamls)
}

func CreateFakePodReadOnlyRootFilesystem(namespace string) {
	yamls := []string{"fakePodRORF1.yml", "fakePodRORF2.yml", "fakePodRORF3.yml", "fakePodRORF4.yml"}
	createHelper(namespace, podsPath, yamls)
}

func CreateFakePodAutomountServiceAccountToken(namespace string) {
	yamls := []string{"fakePodASAT1.yml", "fakePodASAT2.yml", "fakePodASAT3.yml"}
	createHelper(namespace, podsPath, yamls)
}

func CreateFakePodImg(namespace string) {
	yamls := []string{"fakePodImg1.yml", "fakePodImg2.yml"}
	createHelper(namespace, podsPath, yamls)
}
