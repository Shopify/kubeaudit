package fix

import (
	"testing"

	"github.com/Shopify/kubeaudit"
	"github.com/Shopify/kubeaudit/internal/k8s"
	"github.com/Shopify/kubeaudit/k8stypes"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestFix(t *testing.T) {
	cases := []struct {
		testName    string
		pendingFix  kubeaudit.PendingFix
		preFix      func(resource k8stypes.Resource)
		assertFixed func(t *testing.T, resource k8stypes.Resource)
	}{
		{
			testName:   "BySettingPodAnnotation",
			pendingFix: &BySettingPodAnnotation{Key: "mykey", Value: "myvalue"},
			preFix:     func(resource k8stypes.Resource) {},
			assertFixed: func(t *testing.T, resource k8stypes.Resource) {
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
			preFix:     func(resource k8stypes.Resource) {},
			assertFixed: func(t *testing.T, resource k8stypes.Resource) {
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
			preFix: func(resource k8stypes.Resource) {
				k8s.GetPodObjectMeta(resource).SetAnnotations(map[string]string{"mykey": "myvalue"})
			},
			assertFixed: func(t *testing.T, resource k8stypes.Resource) {
				annotations := k8s.GetAnnotations(resource)
				_, ok := annotations["mykey"]
				assert.False(t, ok)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.testName, func(t *testing.T) {
			resource := &k8stypes.PodV1{Spec: v1.PodSpec{}}
			tt.preFix(resource)
			assert.NotEmpty(t, tt.pendingFix.Plan())
			tt.pendingFix.Apply(resource)
			tt.assertFixed(t, resource)
		})
	}
}
