// Copyright 2018 Owkin, inc.
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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultPermission = Permission{
		Public:        false,
		AuthorizedIDs: []string{"foo"},
	}
	defaultPermissions = Permissions{
		Process: defaultPermission,
	}
	defaultOwner = "me"
)

func TestPermissionsCanProcess(t *testing.T) {
	perms := defaultPermissions

	testTable := []struct {
		name           string
		public         bool
		authorizedIDs  []string
		node           string
		expectedAccess bool
	}{
		{"Owner can process", false, []string{}, defaultOwner, true},
		{"Listed node can process", false, []string{"foo"}, "foo", true},
		{"Unlisted node can't process", false, []string{"foo"}, "baz", false},
		{"Everybody can process", true, []string{}, "them", true},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			perms.Process.Public = test.public
			perms.Process.AuthorizedIDs = test.authorizedIDs

			access := perms.CanProcess(defaultOwner, test.node)
			assert.Equal(t, test.expectedAccess, access)
		})
	}
}

func TestPrivInclusion(t *testing.T) {
	testTable := []struct {
		name             string
		includedOpenbar  bool
		includingOpenbar bool
		includedNodes    []string
		includingNodes   []string
		doesInclude      bool
	}{
		{"Full open bar", true, true, []string{}, []string{}, true},
		{"Open bar included but not including is not ok", true, false, []string{}, []string{}, false},
		{"Open bar including but not included is ok", false, true, []string{}, []string{}, true},
		{"One that is both", false, false, []string{"one"}, []string{"one"}, true},
		{"One that is included only", false, false, []string{"one"}, []string{}, false},
		{"One that is including only", false, false, []string{}, []string{"one"}, true},
		{"One is both, two is only including", false, false, []string{"one"}, []string{"one", "two"}, true},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			privIncluded := Permission{Public: test.includedOpenbar, AuthorizedIDs: test.includedNodes}
			privIncluding := Permission{Public: test.includingOpenbar, AuthorizedIDs: test.includingNodes}
			assert.Equal(t, test.doesInclude, privIncluding.include(privIncluded))
		})
	}
}

func TestMergingMechanism(t *testing.T) {
	testTable := []struct {
		name        string
		toMergeINR  bool
		toMergeRU   []string
		expectedINR bool
		expectedRU  []string
	}{
		{"Open bar is not contagious", true, []string{}, false, []string{"foo"}},
		{"The strictest is absolute", false, []string{}, false, []string{}},
		{"Nothing in common", false, []string{"bar"}, false, []string{}},
		{"Only the one in common", false, []string{"foo"}, false, []string{"foo"}},
		{"One in common among of many", false, []string{"foo", "bar", "baz"}, false, []string{"foo"}},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			toMerge := Permission{
				Public:        test.toMergeINR,
				AuthorizedIDs: test.toMergeRU,
			}
			mergedPriv := mergePermissions(defaultPermission, toMerge)
			assert.Equal(t, test.expectedINR, mergedPriv.Public)
			assert.ElementsMatch(t, test.expectedRU, mergedPriv.AuthorizedIDs)
			privMerged := mergePermissions(toMerge, defaultPermission)
			assert.Equal(t, mergedPriv.Public, privMerged.Public, "merging should be transitif")
			assert.ElementsMatch(t, mergedPriv.AuthorizedIDs, privMerged.AuthorizedIDs, "merging should be transitif")

			theSamePriv := mergePermissions(Permission{Public: true}, toMerge)
			assert.Equal(t, toMerge.Public, theSamePriv.Public, "a non restrictive permission should be neutral")
			assert.ElementsMatch(t, toMerge.AuthorizedIDs, theSamePriv.AuthorizedIDs, "a non restrictive permission should be neutral")
			theSamePriv = mergePermissions(toMerge, Permission{Public: true})
			assert.Equal(t, toMerge.Public, theSamePriv.Public, "neutral element should be transitive")
			assert.ElementsMatch(t, toMerge.AuthorizedIDs, theSamePriv.AuthorizedIDs, "neutral element should be transitive")
		})
	}
}
