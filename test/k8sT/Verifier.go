// Copyright 2020 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8sTest

import (
	"fmt"

	. "github.com/cilium/cilium/test/ginkgo-ext"
	"github.com/cilium/cilium/test/helpers"
)

const (
	script = "bpf/verifier-test.sh"
)

var _ = Describe("K8sVerifier", func() {
	var kubectl *helpers.Kubectl

	BeforeAll(func() {
		kubectl = helpers.CreateKubectl(helpers.K8s1VMName(), logger)
		// We don't check the returned error because Cilium could
		// already be removed (e.g., first test to run).
		kubectl.CiliumUninstall("cilium.yaml", map[string]string{})
		ExpectCiliumNotRunning(kubectl)

		kubectl.ExecWithSudo("rm -f /sys/fs/bpf/tc/globals/cilium*").ExpectSuccess("Failed to remove BPF objects")
		cmd := fmt.Sprintf("docker run -v %s/..:/cilium quay.io/cilium/cilium-builder:2020-07-20 make -C /cilium/bpf clean V=0", kubectl.BasePath())
		kubectl.Exec(cmd).ExpectSuccess("Expected cleaning the bpf/ tree to succeed")
	})

	It("Runs the kernel verifier against Cilium's BPF datapath", func() {
		By("Building BPF objects from the tree")
		cmd := fmt.Sprintf("docker run -v %s/..:/cilium quay.io/cilium/cilium-builder:2020-07-20 make -C /cilium/bpf V=0", kubectl.BasePath())
		kubectl.Exec(cmd).ExpectSuccess("Expected compilation of the BPF objects to succeed")
		cmd = fmt.Sprintf("make -C %s/../tools/maptool/", kubectl.BasePath())
		kubectl.Exec(cmd).ExpectSuccess("Expected compilation of maptool to succeed")

		By("Running the verifier test script")
		cmd = kubectl.GetFilePath(script)
		kubectl.ExecWithSudo(cmd).ExpectSuccess("Expected the kernel verifier to pass for BPF programs")
	})
})
