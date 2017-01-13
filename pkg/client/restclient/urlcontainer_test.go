/*
Copyright 2017 The Kubernetes Authors.

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
	"fmt"
	"net/url"
	"testing"
)

func TestURLContainerExclude(t *testing.T) {
	urls := make([]*url.URL, 0, 2)
	for i := 0; i < 2; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://localhost:808%d", i))
		urls = append(urls, u)
	}
	container := NewURLContainer(urls)
	if urls := container.Get(); len(urls) != 2 {
		t.Errorf("Initial length should be equal to provided list of urls: lth %v URLs %v", len(urls), urls)
	}
	container.Exclude(urls[0])
	if urls := container.Get(); len(urls) != 1 {
		t.Errorf("After exclude only single url should be considered as valid: lth %v URLs %v", len(urls), urls)
	}
	container.Exclude(urls[1])
	if urls := container.Get(); len(urls) != 2 {
		t.Errorf("If all urls were excluded we expect that we will considered all urls as valid again: lth %v URLs %v", len(urls), urls)
	}
}
