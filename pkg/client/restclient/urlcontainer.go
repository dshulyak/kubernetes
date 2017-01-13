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
	"net/url"
	"sync"
)

func NewURLContainer(urls []*url.URL) *URLContainer {
	validUrls := make([]*url.URL, len(urls))
	copy(validUrls, urls)
	return &URLContainer{
		urls:      urls,
		validUrls: validUrls,
	}
}

type URLContainer struct {
	m         sync.Mutex
	urls      []*url.URL
	validUrls []*url.URL
}

// Get returns list of valid urls, if valid list of urls is empty - return
// all urls that URLContainer knows about
func (c *URLContainer) Get() []*url.URL {
	c.m.Lock()
	defer c.m.Unlock()
	if len(c.validUrls) != 0 {
		return c.validUrls
	}
	c.validUrls = make([]*url.URL, len(c.urls))
	copy(c.validUrls, c.urls)
	return c.validUrls
}

func (c *URLContainer) Exclude(url *url.URL) {
	c.m.Lock()
	defer c.m.Unlock()
	for i := range c.validUrls {
		if url == c.validUrls[i] {
			c.validUrls = append(c.validUrls[:i], c.validUrls[i+1:]...)
			return
		}
	}
}
