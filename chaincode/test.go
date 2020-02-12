// Run using:
// go build; ./chaincode

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

var key = "b6c1210660962661cd059e0e0608795cc7f61883908aa6df6b0983f42f1d83a0"
var storageAddress = "http://substra-backend.node-2.com/model/82a376bda111134dbfde27dce9301337fcd3dc338b0364f2c99501a364e47d43/file/"
var tt = CompositeTraintuple{
	AssetType:     CompositeTraintupleType,
	AlgoKey:       key,
	ComputePlanID: key,
	Creator:       "MyOrg1MSP",
	Log:           "",
	Rank:          0,
	Status:        "doing",
	Tag:           "",
	Dataset: &Dataset{
		DataManagerKey: key,
		DataSampleKeys: []string{key},
		Worker:         "MyOrg1MSP",
	},
	InHeadModel:  key,
	InTrunkModel: key,
	OutHeadModel: CompositeTraintupleOutHeadModel{
		OutModel: &Hash{
			Hash: key,
		},
		Permissions: Permissions{
			Download: Permission{
				Public:        false,
				AuthorizedIDs: []string{"MyOrg1MSP"},
			},
			Process: Permission{
				Public:        false,
				AuthorizedIDs: []string{"MyOrg1MSP"},
			},
		},
	},
	OutTrunkModel: CompositeTraintupleOutModel{
		OutModel: &HashDress{
			Hash:           key,
			StorageAddress: storageAddress,
		},
		Permissions: Permissions{
			Download: Permission{
				Public:        false,
				AuthorizedIDs: []string{"MyOrg1MSP"},
			},
			Process: Permission{
				Public:        false,
				AuthorizedIDs: []string{"MyOrg1MSP"},
			},
		},
	},
}

type serializer interface {
	name() string
	serialize(object interface{}) ([]byte, error)
	deserialize(serialized []byte, object interface{}) error
}

// JSON
type jsonSerializer struct{}

func (enc jsonSerializer) name() string {
	return "json"
}
func (enc jsonSerializer) serialize(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}
func (enc jsonSerializer) deserialize(serialized []byte, object interface{}) error {
	return json.Unmarshal(serialized, object)
}

// JSON-Gzip
type jsonGzipSerializer struct{}

func (enc jsonGzipSerializer) name() string {
	return "json-gzip"
}
func (enc jsonGzipSerializer) serialize(object interface{}) ([]byte, error) {
	payload, err := json.Marshal(object)
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(payload)
	if err != nil {
		return nil, err
	}
	zw.Close()
	return buf.Bytes(), err
}
func (enc jsonGzipSerializer) deserialize(serialized []byte, object interface{}) error {
	zr, err := gzip.NewReader(bytes.NewReader(serialized))
	if err != nil {
		return err
	}
	payload, err := ioutil.ReadAll(zr)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, object)
}

// gob
type gobSerializer struct{}

func (enc gobSerializer) name() string {
	return "gob"
}
func (enc gobSerializer) serialize(object interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(object)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
func (enc gobSerializer) deserialize(serialized []byte, object interface{}) error {
	d := gob.NewDecoder(bytes.NewReader(serialized))
	return d.Decode(object)
}

// gob-gzip
type gobGzipSerializer struct{}

func (enc gobGzipSerializer) name() string {
	return "gob-gzip"
}
func (enc gobGzipSerializer) serialize(object interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(object)
	if err != nil {
		return nil, err
	}

	payload := b.Bytes()

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(payload)
	if err != nil {
		return nil, err
	}
	zw.Close()
	return buf.Bytes(), err
}

func (enc gobGzipSerializer) deserialize(serialized []byte, object interface{}) error {
	zr, err := gzip.NewReader(bytes.NewReader(serialized))
	if err != nil {
		return err
	}
	payload, err := ioutil.ReadAll(zr)
	if err != nil {
		return err
	}
	d := gob.NewDecoder(bytes.NewReader(payload))
	return d.Decode(object)
}

var serializers = [4]serializer{
	jsonSerializer{},
	jsonGzipSerializer{},
	gobSerializer{},
	gobGzipSerializer{},
}

func main() {
	// validate
	for _, s := range serializers {
		if err := validate(s); err != nil {
			fmt.Printf("%s is NOT valid : %s\n", s.name(), err.Error())
		}
	}

	// run
	for _, s := range serializers {
		run(s)
	}
}

func validate(s serializer) error {
	serialized, err := s.serialize(tt)
	if err != nil {
		return err
	}
	tt2 := Traintuple{}
	err = s.deserialize(serialized, &tt2)
	if err != nil {
		return err
	}
	if tt.AlgoKey != tt2.AlgoKey {
		return fmt.Errorf("Not equal")
	}
	return nil
}

func run(s serializer) {
	start := time.Time{}
	size := 0
	for i := 0; i < 50000; i++ {
		serialized, _ := s.serialize(tt)
		tt2 := CompositeTraintuple{}
		s.deserialize(serialized, &tt2)
		if i == 0 {
			size = len(serialized)
		}
		if i == 10 {
			// enough warmup, now let's measure!
			start = time.Now()
		}
	}
	duration := time.Since(start)
	fmt.Printf("%s: %d ms, %d bytes\n", s.name(), duration.Milliseconds(), size)
}
