package integration

import (
	"bytes"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	userv1client "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	"github.com/openshift/oc/pkg/helpers/tokencmd"
	testutil "github.com/openshift/origin/test/util"
	testserver "github.com/openshift/origin/test/util/server"
)

func TestCLIGetToken(t *testing.T) {
	masterConfig, clusterAdminKubeConfig, err := testserver.StartTestMasterAPI()
	if err != nil {
		t.Fatal(err)
	}
	defer testserver.CleanupMasterEtcd(t, masterConfig)
	clusterAdminClientConfig, err := testutil.GetClusterAdminClientConfig(clusterAdminKubeConfig)
	if err != nil {
		t.Fatal(err)
	}

	anonymousConfig := rest.AnonymousClientConfig(clusterAdminClientConfig)
	reader := bytes.NewBufferString("user\npass")
	accessToken, err := tokencmd.RequestToken(anonymousConfig, reader, "", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(accessToken) == 0 {
		t.Error("Expected accessToken, but did not get one")
	}

	clientConfig := rest.AnonymousClientConfig(clusterAdminClientConfig)
	clientConfig.BearerToken = accessToken

	user, err := userv1client.NewForConfigOrDie(clientConfig).Users().Get("~", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if user.Name != "user" {
		t.Errorf("expected %v, got %v", "user", user.Name)
	}
}
