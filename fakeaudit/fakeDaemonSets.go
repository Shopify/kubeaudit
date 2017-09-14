package fakeaudit

import (
	"path/filepath"
)

var daemonSetPath = filepath.Join(absPath, "fakeaudit", "test", "daemonSets")

func CreateFakeDaemonSetSC(namespace string) {
	fakeDaemonSetClient := getFakeDaemonSetClient(namespace)
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetSC1.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetSC2.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetSC3.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetSC4.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetSC5.yml")))
}

func CreateFakeDaemonSetRunAsNonRoot(namespace string) {
	fakeDaemonSetClient := getFakeDaemonSetClient(namespace)
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRANR1.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRANR2.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRANR3.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRANR4.yml")))
}

func CreateFakeDaemonSetReadOnlyRootFilesystem(namespace string) {
	fakeDaemonSetClient := getFakeDaemonSetClient(namespace)
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRORF1.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRORF2.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRORF3.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetRORF4.yml")))
}

func CreateFakeDaemonSetAutomountServiceAccountToken(namespace string) {
	fakeDaemonSetClient := getFakeDaemonSetClient(namespace)
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetASAT1.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetASAT2.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetASAT3.yml")))
}

func CreateFakeDaemonSetImg(namespace string) {
	fakeDaemonSetClient := getFakeDaemonSetClient(namespace)
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetImg1.yml")))
	fakeDaemonSetClient.Create(getDaemonSet(filepath.Join(daemonSetPath, "fakeDaemonSetImg2.yml")))
}
