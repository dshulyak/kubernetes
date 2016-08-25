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

package restclient

import (
	"path"
	"testing"

	"k8s.io/kubernetes/pkg/api/testapi"
)

func TestValidatesHostParameter(t *testing.T) {
	testCases := []struct {
		Hosts   []string
		APIPath string

		URL string
		Err bool
	}{
		{[]string{"127.0.0.1"}, "", "http://127.0.0.1/" + testapi.Default.GroupVersion().Version, false},
		{[]string{"127.0.0.1:8080"}, "", "http://127.0.0.1:8080/" + testapi.Default.GroupVersion().Version, false},
		{[]string{"foo.bar.com"}, "", "http://foo.bar.com/" + testapi.Default.GroupVersion().Version, false},
		{[]string{"http://host/prefix"}, "", "http://host/prefix/" + testapi.Default.GroupVersion().Version, false},
		{[]string{"http://host"}, "", "http://host/" + testapi.Default.GroupVersion().Version, false},
		{[]string{"http://host"}, "/", "http://host/" + testapi.Default.GroupVersion().Version, false},
		{[]string{"http://host"}, "/other", "http://host/other/" + testapi.Default.GroupVersion().Version, false},
		{[]string{"host/server"}, "", "", true},
	}
	for i, testCase := range testCases {
		basepaths, versionedAPIPath, err := DefaultServerURLsForHosts(testCase.Hosts, testCase.APIPath, *testapi.Default.GroupVersion(), false)
		switch {
		case err == nil && testCase.Err:
			t.Errorf("expected error but was nil")
			continue
		case err != nil && !testCase.Err:
			t.Errorf("unexpected error %v", err)
			continue
		case err != nil:
			continue
		}
		for _, u := range basepaths {
			u.Path = path.Join(u.Path, versionedAPIPath)
			if e, a := testCase.URL, u.String(); e != a {
				t.Errorf("%d: expected host %s, got %s", i, e, a)
				continue
			}
		}
	}
}
