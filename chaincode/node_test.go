package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStub("substra", scc)

	args := append([][]byte{[]byte("registerNode")}, []byte{})

	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "Node created")
	assert.EqualValuesf(t, "{\"id\":\"SampleOrg\"}", string(resp.Payload), "Node created")

	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "Node registered twice")
	assert.EqualValuesf(t, "{\"id\":\"SampleOrg\"}", string(resp.Payload), "Node registered twice")
}

func TestQueryNodes(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	response := mockStub.MockInvoke("43", [][]byte{[]byte("queryNodes")})

	assert.EqualValuesf(t, 200, response.Status, "Node Created")
	assert.EqualValuesf(t, "[{\"id\":\"SampleOrg\"}]", string(response.Payload), "Query nodes")
}
