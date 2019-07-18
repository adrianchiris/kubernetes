/*
Copyright 2015 The Kubernetes Authors.

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

package topologymanager

import (
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager/socketmask"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
)

type preferredPolicy struct{}

var _ Policy = &preferredPolicy{}

// PolicyPreferred policy name.
const PolicyPreferred string = "preferred"

// NewPreferredPolicy returns preferred policy.
func NewPreferredPolicy() Policy {
	return &preferredPolicy{}
}

func (p *preferredPolicy) Name() string {
	return string(PolicyPreferred)
}

func (p *preferredPolicy) CanAdmitPodResult(admit bool) lifecycle.PodAdmitResult {
	return lifecycle.PodAdmitResult{
		Admit: true,
	}
}

func (p *preferredPolicy) Merge(permutation map[string]TopologyHint) TopologyHint {
			// Get the SocketAffinity from each hint in the permutation and see if any
		// of them encode unpreferred allocations.
		preferred := true
		var socketAffinities []socketmask.SocketMask
		for _, hint := range permutation {
			// Only consider hints that have an actual SocketAffinity set.
			if hint.SocketAffinity != nil {
				if !hint.Preferred {
					preferred = false
				}
				socketAffinities = append(socketAffinities, hint.SocketAffinity)
			}
		}

		// Merge the affinities using a bitwise-and operation.
		mergedAffinity, _ := socketmask.NewSocketMask()
		mergedAffinity.Fill()
		mergedAffinity.And(socketAffinities...)

		// Build a mergedHintfrom the merged affinity mask, indicating if an
		// preferred allocation was used to generate the affinity mask or not.
		return TopologyHint{mergedAffinity, preferred}
}