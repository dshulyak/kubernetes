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

package predicates

import (
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/plugin/pkg/scheduler/algorithm"
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
)

type PredicateMetadataFactory struct {
	podLister algorithm.PodLister
}

func NewPredicateMetadataFactory(podLister algorithm.PodLister) algorithm.MetadataProducer {
	factory := &PredicateMetadataFactory{
		podLister,
	}
	return factory.GetMetadata
}

// GetMetadata returns the predicateMetadata used which will be used by various predicates.
func (pfactory *PredicateMetadataFactory) GetMetadata(pod *v1.Pod, nodeNameToInfoMap map[string]*schedulercache.NodeInfo) interface{} {
	// If we cannot compute metadata, just return nil
	if pod == nil {
		return nil
	}
	predicateMetadata := &predicateMetadata{
		pod:           pod,
		podBestEffort: isPodBestEffort(pod),
		podRequest:    GetResourceRequest(pod),
		podPorts:      GetUsedPorts(pod),
	}
	affinity, err := v1.GetAffinityFromPodAnnotations(pod.Annotations)
	if err != nil {
		return nil
	}
	if affinity != nil {
		affinityTerms, err := computeSelectorsAndNamespacesForTerms(pod, getPodAffinityTerms(affinity.PodAffinity))
		if err != nil {
			return nil
		}
		predicateMetadata.affinityTerms = affinityTerms
		antiAffinityTerms, err := computeSelectorsAndNamespacesForTerms(pod, getPodAntiAffinityTerms(affinity.PodAntiAffinity))
		if err != nil {
			return nil
		}
		predicateMetadata.antiAffinityTerms = antiAffinityTerms
		matchingTerms, err := getMatchingAntiAffinityTerms(pod, nodeNameToInfoMap, antiAffinityTerms)
		if err != nil {
			return nil
		}
		predicateMetadata.matchingAntiAffinityTerms = matchingTerms
	}
	for predicateName, precomputeFunc := range predicatePrecomputations {
		glog.V(10).Info("Precompute: %v", predicateName)
		precomputeFunc(predicateMetadata)
	}
	return predicateMetadata
}
