package cmd

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

type Result struct {
	Err            int
	Occurrences    []Occurrence
	Namespace      string
	Name           string
	CapsAdded      []Capability
	ImageName      string
	CapsDropped    []Capability
	CapsNotDropped []Capability
	KubeType       string
	DSA            string
	SA             string
	Token          *bool
	ImageTag       string
	CPULimitActual string
	CPULimitMax    string
	MEMLimitActual string
	MEMLimitMax    string
}

func (res Result) Print() {
	for _, occ := range res.Occurrences {
		if occ.kind <= KubeauditLogLevels[rootConfig.verbose] {
			logger := log.WithFields(createFields(res, occ.id))
			switch occ.kind {
			case Debug:
				logger.Debug(occ.message)
			case Info:
				logger.Info(occ.message)
			case Warn:
				logger.Warn(occ.message)
			case Error:
				logger.Error(occ.message)
			}
		}
	}
}

func createFields(res Result, err int) (fields log.Fields) {
	fields = log.Fields{}
	v := reflect.ValueOf(res)
	for _, member := range shouldLog(err) {
		value := v.FieldByName(member)
		if value.IsValid() && value.Interface() != nil && value.Interface() != "" {
			fields[member] = value.Interface()
		}
	}
	return
}

func shouldLog(err int) (members []string) {
	members = []string{"Name", "Namespace", "KubeType"}
	switch err {
	case ErrorCapabilitiesAdded:
		members = append(members, "CapsAdded")
	case ErrorCapabilitiesSomeDropped:
		members = append(members, "CapsNotDropped")
	case ErrorServiceAccountTokenDeprecated:
		members = append(members, "DSA")
		members = append(members, "SA")
	case InfoImageCorrect:
	case ErrorImageTagMissing:
	case ErrorImageTagIncorrect:
		members = append(members, "ImageTag")
		members = append(members, "ImageName")
	case ErrorResourcesLimitsCpuExceeded:
		members = append(members, "CPULimitActual")
		members = append(members, "CPULimitMax")
	case ErrorResourcesLimitsMemoryExceeded:
		members = append(members, "MEMLimitActual")
		members = append(members, "MEMLimitMax")
	}
	return
}
