package e2e

import (
	"fmt"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = framework.KubeDescribe("memhog [Conformance]", func() {
	f := framework.NewDefaultFramework("memhog")
	It("memhog pod should be evicted", func() {
		besteffort := framework.CreateMemhogPod(f, "besteffort-", "besteffort", api.ResourceRequirements{})
		Eventually(func() error {
			glog.Infof("waiting for memhog pod to be evicted")
			best, err := f.Client.Pods(f.Namespace.Name).Get(besteffort.Name)
			framework.ExpectNoError(err)
			if best.Status.Phase == api.PodFailed {
				return nil
			}
			return fmt.Errorf("besteffort is not evicted.")
		}, 20*time.Minute, 5*time.Second).Should(BeNil())
	})
})
