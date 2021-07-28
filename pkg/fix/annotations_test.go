package fix

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/pkg/k8s"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestFix(t *testing.T) {
	cases := []struct {
		testName    string
		pendingFix  kubeaudit.PendingFix
		preFix      func(resource k8s.Resource)
		assertFixed func(t *testing.T, resource k8s.Resource)
	}{
		{
			testName:   "BySettingPodAnnotation",
			pendingFix: &BySettingPodAnnotation{Key: "mykey", Value: "myvalue"},
			preFix:     func(resource k8s.Resource) {},
			assertFixed: func(t *testing.T, resource k8s.Resource) {
				annotations := k8s.GetAnnotations(resource)
				assert.NotNil(t, annotations)
				val, ok := annotations["mykey"]
				assert.True(t, ok)
				assert.Equal(t, "myvalue", val)
			},
		},
		{
			testName:   "ByAddingPodAnnotation",
			pendingFix: &ByAddingPodAnnotation{Key: "mykey", Value: "myvalue"},
			preFix:     func(resource k8s.Resource) {},
			assertFixed: func(t *testing.T, resource k8s.Resource) {
				annotations := k8s.GetAnnotations(resource)
				assert.NotNil(t, annotations)
				val, ok := annotations["mykey"]
				assert.True(t, ok)
				assert.Equal(t, "myvalue", val)
			},
		},
		{
			testName:   "ByRemovingPodAnnotation",
			pendingFix: &ByRemovingPodAnnotation{Key: "mykey"},
			preFix: func(resource k8s.Resource) {
				k8s.GetPodObjectMeta(resource).SetAnnotations(map[string]string{"mykey": "myvalue"})
			},
			assertFixed: func(t *testing.T, resource k8s.Resource) {
				annotations := k8s.GetAnnotations(resource)
				_, ok := annotations["mykey"]
				assert.False(t, ok)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			resource := &k8s.PodV1{Spec: v1.PodSpec{}}
			tc.preFix(resource)
			assert.NotEmpty(t, tc.pendingFix.Plan())
			tc.pendingFix.Apply(resource)
			tc.assertFixed(t, resource)
		})
	}
}
