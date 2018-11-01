package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/version"
	fakediscovery "k8s.io/client-go/discovery/fake"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
)

type TestK8sClientInCluster struct{}

func (kc TestK8sClientInCluster) InClusterConfig() (*rest.Config, error) {
	return &rest.Config{}, nil
}

type TestK8sClientNotInCluster struct{}

func (kc TestK8sClientNotInCluster) InClusterConfig() (*rest.Config, error) {
	return nil, fmt.Errorf("unable to load in-cluster configuration, mock error")
}
func TestKubeClientConfig(t *testing.T) {
	type k8sCfgTest struct {
		env    map[string]string
		client Client
		cfg    rootFlags
		log    string
		msg    string
	}

	cfgWithNonExistentConfigFile := rootFlags{kubeConfig: "/notarealfile"}
	cfgWithNoConfigFile := rootFlags{}
	envWithHome := map[string]string{"HOME": "/Users/example"}
	envWithoutHome := map[string]string{"HOME": ""}

	k8sCfgTests := []k8sCfgTest{
		k8sCfgTest{env: envWithHome, client: TestK8sClientInCluster{}, cfg: cfgWithNonExistentConfigFile,
			log: "Unable to load kubeconfig. Could not open file",
			msg: "should have tried to open specified kubeconfig file when running in cluster",
		},
		k8sCfgTest{env: envWithHome, client: TestK8sClientNotInCluster{}, cfg: cfgWithNonExistentConfigFile,
			log: "Unable to load kubeconfig. Could not open file",
			msg: "should have tried to open specified kubeconfig file when not running in cluster",
		},
		k8sCfgTest{env: envWithHome, client: TestK8sClientInCluster{}, cfg: cfgWithNoConfigFile,
			log: "Running inside cluster, using the cluster config",
			msg: "should have tried to use cluster mode when running in cluster and no kubeconfig file specified",
		},
		k8sCfgTest{env: envWithHome, client: TestK8sClientNotInCluster{}, cfg: cfgWithNoConfigFile,
			log: "Not running inside cluster, using local config",
			msg: "should have tried to use default kubeconfig when not running in cluster and no kubeconfig file specified",
		},
		k8sCfgTest{env: envWithoutHome, client: TestK8sClientNotInCluster{}, cfg: cfgWithNoConfigFile,
			log: "Unable to load kubeconfig. No config file specified and $HOME not found",
			msg: "should have failed to find default kubeconfig when $HOME does not exist",
		},
	}

	oldLogOut := log.StandardLogger().Out

	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	for _, tt := range k8sCfgTests {
		oldEnv := make(map[string]string)
		changeEnv(oldEnv, tt.env)

		rootConfig = tt.cfg
		buf.Reset()
		kubeClientConfig(tt.client)
		logs := buf.String()
		assert.Contains(t, logs, tt.log, tt.msg)

		resetEnv(oldEnv, tt.env)
	}

	log.SetOutput(oldLogOut)
}

func changeEnv(oldEnv map[string]string, newEnv map[string]string) {
	for k, v := range newEnv {
		if oldVal, exists := os.LookupEnv(k); exists {
			oldEnv[k] = oldVal
		}
		os.Setenv(k, v)
	}
}

func resetEnv(oldEnv map[string]string, newEnv map[string]string) {
	for k := range newEnv {
		v, exists := oldEnv[k]
		if !exists {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}
}

func TestKubeClientConfigLocal(t *testing.T) {
	rootConfig = rootFlags{
		kubeConfig: "/notarealfile",
	}

	_, err := kubeClientConfigLocal()
	assert.Equal(t, ErrNoReadableKubeConfig, err,
		"kubeClientConfigLocal did not return expected error when kubeconfig file doesn't exist")
}

func TestGetKubernetesVersion(t *testing.T) {
	client := fakeclientset.NewSimpleClientset()
	fakeDiscovery, ok := client.Discovery().(*fakediscovery.FakeDiscovery)
	if !ok {
		t.Fatalf("couldn't mock server version")
	}

	fakeDiscovery.FakedServerVersion = &version.Info{
		Major:     "0",
		Minor:     "0",
		GitCommit: "0000",
		Platform:  "ACME 8-bit",
	}

	r, err := getKubernetesVersion(client)
	assert.Nil(t, err)
	assert.EqualValues(t, *fakeDiscovery.FakedServerVersion, *r)
}
