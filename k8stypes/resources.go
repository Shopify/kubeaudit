package k8stypes

var podTemplateSpec = PodTemplateSpecV1{
	ObjectMeta: ObjectMetaV1{},
	Spec:       PodSpecV1{},
}

// NewDeployment creates a new Deployment resource
func NewDeployment() *DeploymentV1 {
	deployment := &DeploymentV1{
		ObjectMeta: ObjectMetaV1{},
		Spec: DeploymentSpecV1{
			Template: podTemplateSpec,
		},
	}

	deployment.Kind = "Deployment"
	deployment.APIVersion = "apps/v1"
	return deployment
}

// NewPod creates a new Pod resource
func NewPod() *PodV1 {
	pod := &PodV1{
		ObjectMeta: ObjectMetaV1{},
		Spec:       PodSpecV1{},
	}

	pod.Kind = "Pod"
	pod.APIVersion = "v1"
	return pod
}

// NewNamespace creates a new Namespace resource
func NewNamespace() *NamespaceV1 {
	namespace := &NamespaceV1{
		ObjectMeta: ObjectMetaV1{},
		Spec:       NamespaceSpecV1{},
	}

	namespace.Kind = "Namespace"
	namespace.APIVersion = "v1"
	return namespace
}

// NewDaemonSet creates a new DaemonSet resource
func NewDaemonSet() *DaemonSetV1 {
	daemonset := &DaemonSetV1{
		ObjectMeta: ObjectMetaV1{},
		Spec: DaemonSetSpecV1{
			Template: podTemplateSpec,
		},
	}

	daemonset.Kind = "Daemonset"
	daemonset.APIVersion = "apps/v1"
	return daemonset
}

// NewReplicationController creates a new ReplicationController resource
func NewReplicationController() *ReplicationControllerV1 {
	controller := &ReplicationControllerV1{
		ObjectMeta: ObjectMetaV1{},
		Spec: ReplicationControllerSpecV1{
			Template: podTemplateSpec.DeepCopy(),
		},
	}

	controller.Kind = "ReplicationController"
	controller.APIVersion = "v1"
	return controller
}

// NewStatefulSet creates a new StatefulSet resource
func NewStatefulSet() *StatefulSetV1 {
	controller := &StatefulSetV1{
		ObjectMeta: ObjectMetaV1{},
		Spec: StatefulSetSpecV1{
			Template: podTemplateSpec,
		},
	}

	controller.Kind = "StatefulSet"
	controller.APIVersion = "apps/v1"
	return controller
}

// NewNetworkPolicy creates a new NetworkPolicy resource
func NewNetworkPolicy() *NetworkPolicyV1 {
	networkPolicy := &NetworkPolicyV1{
		ObjectMeta: ObjectMetaV1{},
		Spec:       NetworkPolicySpecV1{},
	}

	networkPolicy.Kind = "NetworkPolicy"
	networkPolicy.APIVersion = "networking.k8s.io/v1"
	return networkPolicy
}

// NewPodTemplate creates a new PodTemplate resource
func NewPodTemplate() *PodTemplateV1 {
	podTemplate := &PodTemplateV1{
		ObjectMeta: ObjectMetaV1{},
		Template:   podTemplateSpec,
	}

	podTemplate.Kind = "PodTemplate"
	podTemplate.APIVersion = "core/v1"
	return podTemplate
}

// NewCronJob creates a new CronJob resource
func NewCronJob() *CronJobV1Beta1 {
	cronJob := &CronJobV1Beta1{
		Spec: CronJobSpecV1Beta1{
			JobTemplate: JobTemplateSpecV1Beta1{
				Spec: JobSpecV1{
					Template: podTemplateSpec,
				},
			},
		},
	}

	cronJob.Kind = "CronJob"
	cronJob.APIVersion = "batch/v1beta1"
	return cronJob
}
