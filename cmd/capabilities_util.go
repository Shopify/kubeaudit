package cmd

import "reflect"

func dropCapFromConfigList(capList *KubeauditConfigCapabilities) (dropCapSet CapSet) {
	var configCapabilityValue reflect.Value
	var r reflect.Value
	dropCapSet = make(CapSet)
	r = reflect.ValueOf(capList)
	configCapabilityValue = reflect.Indirect(r).FieldByName("SetPCAP")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["SETPCAP"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("MKNOD")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["MKNOD"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("AuditWrite")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["AUDIT_WRITE"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("Chown")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["CHOWN"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("NetRaw")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["NET_RAW"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("DacOverride")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["DAC_OVERRIDE"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("FOWNER")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["FOWNER"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("FSetID")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["FSETID"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("Kill")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["KILL"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("SetGID")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["SETGID"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("SetUID")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["SETUID"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("NetBindService")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["NET_BIND_SERVICE"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("SYSChroot")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["SYS_CHROOT"] = true
	}

	configCapabilityValue = reflect.Indirect(r).FieldByName("SetFCAP")
	if configCapabilityValue.String() == "drop" {
		dropCapSet["SETFCAP"] = true
	}

	return dropCapSet
}
