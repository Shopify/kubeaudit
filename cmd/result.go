package cmd

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

type Result struct {
	CPULimitActual string
	CPULimitMax    string
	DSA            string
	Err            int
	ImageName      string
	ImageTag       string
	KubeType       string
	Labels         map[string]string
	MEMLimitActual string
	MEMLimitMax    string
	Name           string
	Namespace      string
	Occurrences    []Occurrence
	SA             string
	Token          *bool
}

func (res Result) Print() {
	for _, occ := range res.Occurrences {
		if occ.kind <= KubeauditLogLevels[rootConfig.verbose] {
			logger := log.WithFields(createFields(res, occ))
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

func createFields(res Result, occ Occurrence) (fields log.Fields) {
	fields = log.Fields{}
	v := reflect.ValueOf(res)
	for _, member := range shouldLog(occ.id) {
		value := v.FieldByName(member)
		if value.IsValid() && value.Interface() != nil && value.Interface() != "" {
			fields[member] = value.Interface()
		}
	}
	for k, v := range occ.metadata {
		fields[k] = v
	}
	return
}

func shouldLog(err int) (members []string) {
	members = []string{"Name", "Namespace", "KubeType"}
	switch err {
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
