package fakeaudit

import (
	"path/filepath"
)

var daemonSetPath = filepath.Join(absPath, "fakeaudit", "test", "daemonSets")

func CreateFakeDaemonSetSC(namespace string) {
	yamls := []string{"fakeDaemonSetSC1.yml", "fakeDaemonSetSC2.yml", "fakeDaemonSetSC3.yml", "fakeDaemonSetSC4.yml", "fakeDaemonSetSC5.yml"}
	createHelper(namespace, daemonSetPath, yamls)
}

func CreateFakeDaemonSetRunAsNonRoot(namespace string) {
	yamls := []string{"fakeDaemonSetRANR1.yml", "fakeDaemonSetRANR2.yml", "fakeDaemonSetRANR2.yml", "fakeDaemonSetRANR3.yml", "fakeDaemonSetRANR4.yml"}
	createHelper(namespace, daemonSetPath, yamls)
}

func CreateFakeDaemonSetPrivileged(namespace string) {
	yamls := []string{"fakeDaemonSetPrivileged1.yml", "fakeDaemonSetPrivileged2.yml", "fakeDaemonSetPrivileged3.yml", "fakeDaemonSetPrivileged4.yml"}
	createHelper(namespace, daemonSetPath, yamls)
}

func CreateFakeDaemonSetReadOnlyRootFilesystem(namespace string) {
	yamls := []string{"fakeDaemonSetRORF1.yml", "fakeDaemonSetRORF2.yml", "fakeDaemonSetRORF3.yml", "fakeDaemonSetRORF4.yml"}
	createHelper(namespace, daemonSetPath, yamls)
}

func CreateFakeDaemonSetAutomountServiceAccountToken(namespace string) {
	yamls := []string{"fakeDaemonSetASAT1.yml", "fakeDaemonSetASAT2.yml", "fakeDaemonSetASAT3.yml"}
	createHelper(namespace, daemonSetPath, yamls)
}

func CreateFakeDaemonSetImg(namespace string) {
	yamls := []string{"fakeDaemonSetImg1.yml", "fakeDaemonSetImg2.yml"}
	createHelper(namespace, daemonSetPath, yamls)
}
