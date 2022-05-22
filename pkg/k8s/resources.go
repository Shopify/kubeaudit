package k8s

var podTemplateSpec = PodTemplateSpecV1{
	ObjectMeta: ObjectMetaV1{},
	Spec:       PodSpecV1{},
}

// NewDeployment creates a new Deployment resource
func NewDeployment() *DeploymentV1 {
	return &DeploymentV1{
		TypeMeta: TypeMetaV1{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec: DeploymentSpecV1{
			Template: podTemplateSpec,
		},
	}
}

// NewPod creates a new Pod resource
func NewPod() *PodV1 {
	return &PodV1{
		TypeMeta: TypeMetaV1{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec:       PodSpecV1{},
	}
}

// NewNamespace creates a new Namespace resource
func NewNamespace() *NamespaceV1 {
	return &NamespaceV1{
		TypeMeta: TypeMetaV1{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec:       NamespaceSpecV1{},
	}
}

// NewDaemonSet creates a new DaemonSet resource
func NewDaemonSet() *DaemonSetV1 {
	return &DaemonSetV1{
		TypeMeta: TypeMetaV1{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec: DaemonSetSpecV1{
			Template: podTemplateSpec,
		},
	}
}

// NewReplicationController creates a new ReplicationController resource
func NewReplicationController() *ReplicationControllerV1 {
	return &ReplicationControllerV1{
		TypeMeta: TypeMetaV1{
			Kind:       "ReplicationController",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec: ReplicationControllerSpecV1{
			Template: podTemplateSpec.DeepCopy(),
		},
	}
}

// NewStatefulSet creates a new StatefulSet resource
func NewStatefulSet() *StatefulSetV1 {
	return &StatefulSetV1{
		TypeMeta: TypeMetaV1{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec: StatefulSetSpecV1{
			Template: podTemplateSpec,
		},
	}
}

// NewNetworkPolicy creates a new NetworkPolicy resource
func NewNetworkPolicy() *NetworkPolicyV1 {
	return &NetworkPolicyV1{
		TypeMeta: TypeMetaV1{
			Kind:       "NetworkPolicy",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec:       NetworkPolicySpecV1{},
	}
}

// NewPodTemplate creates a new PodTemplate resource
func NewPodTemplate() *PodTemplateV1 {
	return &PodTemplateV1{
		TypeMeta: TypeMetaV1{
			Kind:       "PodTemplate",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Template:   podTemplateSpec,
	}
}

// NewCronJob creates a new CronJob resource
func NewCronJob() *CronJobV1Beta1 {
	return &CronJobV1Beta1{
		TypeMeta: TypeMetaV1{
			Kind:       "CronJob",
			APIVersion: "batch/v1beta1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec: CronJobSpecV1Beta1{
			JobTemplate: JobTemplateSpecV1Beta1{
				Spec: JobSpecV1{
					Template: podTemplateSpec,
				},
			},
		},
	}
}

// NewServiceAccount creates a new ServiceAccount resource
func NewServiceAccount() *ServiceAccountV1 {
	return &ServiceAccountV1{
		TypeMeta: TypeMetaV1{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMetaV1{},
	}
}

// NewService creates a new Service resource
func NewService() *ServiceV1 {
	return &ServiceV1{
		TypeMeta: TypeMetaV1{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMetaV1{},
	}
}

// NewJob creates a new Job resource
func NewJob() *JobV1 {
	return &JobV1{
		TypeMeta: TypeMetaV1{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: ObjectMetaV1{},
		Spec: JobSpecV1{
			Template: podTemplateSpec,
		},
	}
}
