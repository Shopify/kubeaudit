package cmd

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Result stores information about a Kubernetes resource, including all audit results (Occurrences) related to that
// resource.
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
	FileName       string
	Namespace      string
	Occurrences    []Occurrence
	SA             string
	Token          *bool
}

// Print logs all audit results to their respective log levels.
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
	fields["Container"] = occ.container
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
	case ErrorResourcesLimitsCPUExceeded:
		members = append(members, "CPULimitActual")
		members = append(members, "CPULimitMax")
	case ErrorResourcesLimitsMemoryExceeded:
		members = append(members, "MEMLimitActual")
		members = append(members, "MEMLimitMax")
	}
	return
}

func (res *Result) allowedCaps() (allowed map[CapabilityV1]string) {
	allowed = make(map[CapabilityV1]string)
	for k, v := range res.Labels {
		if strings.Contains(k, "audit.kubernetes.io/allow-capability-") {
			capName := strings.Replace(strings.ToUpper(strings.TrimPrefix(k, "audit.kubernetes.io/allow-capability-")), "-", "_", -1)
			allowed[CapabilityV1(capName)] = v
		}
	}
	return
}
