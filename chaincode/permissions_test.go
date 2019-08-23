package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultPrivilege = Privilege{
		IsRestricted:    true,
		AuthorizedNodes: []string{"foo"},
	}
	defaultPermissions = Permissions{
		Process: defaultPrivilege,
	}
	defaultOwner = "me"
)

func TestPermissionsCanProcess(t *testing.T) {
	perms := defaultPermissions

	testTable := []struct {
		name            string
		isRestricted    bool
		authorizedNodes []string
		node            string
		expectedAccess  bool
	}{
		{"Owner can process", true, []string{}, defaultOwner, true},
		{"Listed node can process", true, []string{"foo"}, "foo", true},
		{"Unlisted node can't process", true, []string{"foo"}, "baz", false},
		{"Everybody can process", false, []string{}, "them", true},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			perms.Process.IsRestricted = test.isRestricted
			perms.Process.AuthorizedNodes = test.authorizedNodes

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
		{"Full open bar", false, false, []string{}, []string{}, true},
		{"Open bar included but not including is not ok", false, true, []string{}, []string{}, false},
		{"Open bar including but not included is ok", true, false, []string{}, []string{}, true},
		{"One that is both", true, true, []string{"one"}, []string{"one"}, true},
		{"One that is included only", true, true, []string{"one"}, []string{}, false},
		{"One that is including only", true, true, []string{}, []string{"one"}, true},
		{"One is both, two is only including", true, true, []string{"one"}, []string{"one", "two"}, true},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			privIncluded := Privilege{IsRestricted: test.includedOpenbar, AuthorizedNodes: test.includedNodes}
			privIncluding := Privilege{IsRestricted: test.includingOpenbar, AuthorizedNodes: test.includingNodes}
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
		{"Open bar is not contagious", false, []string{}, true, []string{"foo"}},
		{"The strictest is absolute", true, []string{}, true, []string{}},
		{"Nothing in common", true, []string{"bar"}, true, []string{}},
		{"Only the one in common", true, []string{"foo"}, true, []string{"foo"}},
		{"One in common among of many", true, []string{"foo", "bar", "baz"}, true, []string{"foo"}},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			toMerge := Privilege{
				IsRestricted:    test.toMergeINR,
				AuthorizedNodes: test.toMergeRU,
			}
			mergedPriv := mergePrivileges(defaultPrivilege, toMerge)
			assert.Equal(t, test.expectedINR, mergedPriv.IsRestricted)
			assert.ElementsMatch(t, test.expectedRU, mergedPriv.AuthorizedNodes)
			privMerged := mergePrivileges(toMerge, defaultPrivilege)
			assert.Equal(t, mergedPriv.IsRestricted, privMerged.IsRestricted, "merging should be transitif")
			assert.ElementsMatch(t, mergedPriv.AuthorizedNodes, privMerged.AuthorizedNodes, "merging should be transitif")

			theSamePriv := mergePrivileges(Privilege{IsRestricted: false}, toMerge)
			assert.Equal(t, toMerge.IsRestricted, theSamePriv.IsRestricted, "a non restrictive privilege should be neutral")
			assert.ElementsMatch(t, toMerge.AuthorizedNodes, theSamePriv.AuthorizedNodes, "a non restrictive privilege should be neutral")
			theSamePriv = mergePrivileges(toMerge, Privilege{IsRestricted: false})
			assert.Equal(t, toMerge.IsRestricted, theSamePriv.IsRestricted, "neutral element should be transitive")
			assert.ElementsMatch(t, toMerge.AuthorizedNodes, theSamePriv.AuthorizedNodes, "neutral element should be transitive")
		})
	}
}
