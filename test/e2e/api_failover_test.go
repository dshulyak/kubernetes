/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	// nodeUpdateTimeout is the time to wait for node to change its state.
	nodeUpdateTimeout = 1 * time.Minute
	// nodeUpdatePollInterval is the interval to check ndoe state.
	nodeUpdatePollInterval = 5 * time.Second
)

func waitNodes(f *framework.Framework, condition bool) {
	nodesReady := false
	Eventually(func() error {
		nodes, err := f.Client.Nodes().List(api.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		if len(nodes.Items) == 0 {
			return fmt.Errorf("empty node list: %+v", nodes)
		}
		for _, node := range nodes.Items {
			nodesReady = api.IsNodeReady(&node)
		}
		if condition != nodesReady {
			return fmt.Errorf("Nodes are not in expected condition")
		}
		return nil
	}, nodeUpdateTimeout, nodeUpdatePollInterval).Should(Succeed())
	return
}

var _ = framework.KubeDescribe("Nodes api failover [Disruptive]", func() {

	const (
		// nodeUpdateTimeout is the time to wait for node to change its state.
		nodeUpdateTimeout = 1 * time.Minute
		// nodeUpdatePollInterval is the interval to check ndoe state.
		nodeUpdatePollInterval = 5 * time.Second
	)

	f := framework.NewDefaultFramework("nodes-failover")

	ports := []int{6444, 443}
	It("should not work if both ports are unreachable", func() {
		for _, port := range ports {
			framework.MasterExec(fmt.Sprintf("sudo iptables -A INPUT -p tcp --destination-port %d -j DROP", port))
		}
		waitNodes(f, false)
		for _, port := range ports {
			framework.MasterExec(fmt.Sprintf("sudo iptables -D INPUT -p tcp --destination-port %d -j DROP", port))
		}
	})

	It("should work only if first port is unreachable", func() {
		port := ports[0]
		framework.MasterExec(fmt.Sprintf("sudo iptables -A INPUT -p tcp --destination-port %d -j DROP", port))
		waitNodes(f, true)
		framework.MasterExec(fmt.Sprintf("sudo iptables -D INPUT -p tcp --destination-port %d -j DROP", port))
	})
})
