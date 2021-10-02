// Copyright 2019 The GCR Cleaner Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Test gcrcleaner
package gcrcleaner

import (
	"testing"
	"time"

	gcrgoogle "github.com/google/go-containerregistry/pkg/v1/google"
)

var future, _ = time.Parse("2006-01-02", "2030-01-01")

// 5 years in the past
var past, _ = time.Parse("2006-01-02", "2020-01-01")

var tests = []struct {
	goal     string
	manifest gcrgoogle.ManifestInfo
	since    time.Time
	expect   bool
}{
	{
		goal:  "do not delete manifests with tags",
		since: future,
		manifest: gcrgoogle.ManifestInfo{
			Size:      0,
			MediaType: "",
			Created:   past,
			Uploaded:  past,
			Tags:      []string{"dac-contenthub-latest", "dac-contenthub-latest-2020-10-01", "dac-contenthub-latest-2020-01-01"},
		},
		expect: false,
	},
	{
		since: future,
		goal:  "Delete if we do not have a latest tag and the image was created in the past",
		manifest: gcrgoogle.ManifestInfo{
			Size:      0,
			MediaType: "",
			Created:   past,
			Uploaded:  past,
			Tags:      []string{"dac-contenthub-latest-2029-10-01", "dac-contenthub-latest-2020-01-01"},
		},
		expect: true,
	},
	{
		since: past,
		goal:  "Do not delete if we have not reached the expiration",
		manifest: gcrgoogle.ManifestInfo{
			Size:      0,
			MediaType: "",
			Created:   future,
			Uploaded:  future,
			Tags:      []string{"dac-contenthub-latest-2020-10-01", "dac-contenthub-latest-2020-01-01"},
		},
		expect: false,
	},
	{
		since: future,
		goal:  "Delete if we do not have any tags",
		manifest: gcrgoogle.ManifestInfo{
			Size:      0,
			MediaType: "",
			Created:   past,
			Uploaded:  past,
			Tags:      []string{},
		},
		expect: true,
	},
}

// shouldDelete returns true if the manifest has no tags or allows deletion of tagged images
// and is before the requested time.
func TestShouldDelete(t *testing.T) {
	var c = &Cleaner{}
	for _, test := range tests {
		res := c.shouldDelete(test.manifest, test.since, false, nil)
		if res != test.expect {
			t.Errorf("goal: %s \n - since %v shouldDelete(%v, %v) = %v, expected %v", test.goal, test.since, test.manifest, test.since, res, test.expect)
		}
	}
}
