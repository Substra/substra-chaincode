### Examples

#### ------------ Add Node ------------
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerNode"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "id": "SampleOrg"
}
```
#### ------------ Add DataManager ------------
Smart contract: `registerDataManager`

##### JSON Inputs:
```go
{
 "name": string (required,gte=1,lte=100),
 "openerHash": string (required,len=64,hexadecimal),
 "openerStorageAddress": string (required,url),
 "type": string (required,gte=1,lte=30),
 "descriptionHash": string (required,len=64,hexadecimal),
 "descriptionStorageAddress": string (required,url),
 "objectiveKey": string (omitempty),
 "permissions": (required){
   "process": (required){
     "public": bool (required),
     "authorizedIDs": [string] (required),
   },
 },
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataManager","{\"name\":\"liver slide\",\"openerHash\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"openerStorageAddress\":\"https://toto/dataManager/42234/opener\",\"type\":\"images\",\"descriptionHash\":\"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee\",\"descriptionStorageAddress\":\"https://toto/dataManager/42234/description\",\"objectiveKey\":\"\",\"permissions\":{\"process\":{\"public\":true,\"authorizedIDs\":[]}}}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
}
```
#### ------------ Query DataManager From key ------------
Smart contract: `queryDataManager`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryDataManager","{\"key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "description": {
  "hash": "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee",
  "storageAddress": "https://toto/dataManager/42234/description"
 },
 "key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
 "name": "liver slide",
 "objectiveKey": "",
 "opener": {
  "hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "storageAddress": "https://toto/dataManager/42234/opener"
 },
 "owner": "SampleOrg",
 "permissions": {
  "process": {
   "authorizedIDs": [],
   "public": true
  }
 },
 "type": "images"
}
```
#### ------------ Add test DataSample ------------
Smart contract: `registerDataSample`

##### JSON Inputs:
```go
{
 "hashes": [string] (required,dive,len=64,hexadecimal),
 "dataManagerKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"hashes\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"dataManagerKeys\":[\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"testOnly\":\"true\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "keys": [
  "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
 ]
}
```
#### ------------ Add Objective ------------
Smart contract: `registerObjective`

##### JSON Inputs:
```go
{
 "name": string (required,gte=1,lte=100),
 "descriptionHash": string (required,len=64,hexadecimal),
 "descriptionStorageAddress": string (required,url),
 "metricsName": string (required,gte=1,lte=100),
 "metricsHash": string (required,len=64,hexadecimal),
 "metricsStorageAddress": string (required,url),
 "testDataset": (omitempty){
   "dataManagerKey": string (omitempty,len=64,hexadecimal),
   "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
 },
 "permissions": (required){
   "process": (required){
     "public": bool (required),
     "authorizedIDs": [string] (required),
   },
 },
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerObjective","{\"name\":\"MSI classification\",\"descriptionHash\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"descriptionStorageAddress\":\"https://toto/objective/222/description\",\"metricsName\":\"accuracy\",\"metricsHash\":\"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"metricsStorageAddress\":\"https://toto/objective/222/metrics\",\"testDataset\":{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"]},\"permissions\":{\"process\":{\"public\":true,\"authorizedIDs\":[]}}}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
}
```
#### ------------ Add Algo ------------
Smart contract: `registerAlgo`

##### JSON Inputs:
```go
{
 "name": string (required,gte=1,lte=100),
 "hash": string (required,len=64,hexadecimal),
 "storageAddress": string (required,url),
 "descriptionHash": string (required,len=64,hexadecimal),
 "descriptionStorageAddress": string (required,url),
 "permissions": (required){
   "process": (required){
     "public": bool (required),
     "authorizedIDs": [string] (required),
   },
 },
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo","{\"name\":\"hog + svm\",\"hash\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"storageAddress\":\"https://toto/algo/222/algo\",\"descriptionHash\":\"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca\",\"descriptionStorageAddress\":\"https://toto/algo/222/description\",\"permissions\":{\"process\":{\"public\":true,\"authorizedIDs\":[]}}}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
}
```
#### ------------ Add Train DataSample ------------
Smart contract: `registerDataSample`

##### JSON Inputs:
```go
{
 "hashes": [string] (required,dive,len=64,hexadecimal),
 "dataManagerKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"hashes\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"dataManagerKeys\":[\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"testOnly\":\"false\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "keys": [
  "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
 ]
}
```
#### ------------ Query DataManagers ------------
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryDataManagers"]}' -C myc
```
##### Command output:
```json
[
 {
  "description": {
   "hash": "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee",
   "storageAddress": "https://toto/dataManager/42234/description"
  },
  "key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "liver slide",
  "objectiveKey": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "opener": {
   "hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "storageAddress": "https://toto/dataManager/42234/opener"
  },
  "owner": "SampleOrg",
  "permissions": {
   "process": {
    "authorizedIDs": [],
    "public": true
   }
  },
  "type": "images"
 }
]
```
#### ------------ Query DataSamples ------------
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryDataSamples"]}' -C myc
```
##### Command output:
```json
[
 {
  "dataManagerKeys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "dataManagerKeys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "dataManagerKeys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "dataManagerKeys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 }
]
```
#### ------------ Query Objectives ------------
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryObjectives"]}' -C myc
```
##### Command output:
```json
[
 {
  "description": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/description"
  },
  "key": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "name": "accuracy",
   "storageAddress": "https://toto/objective/222/metrics"
  },
  "name": "MSI classification",
  "owner": "SampleOrg",
  "permissions": {
   "process": {
    "authorizedIDs": [],
    "public": true
   }
  },
  "testDataset": {
   "dataManagerKey": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "dataSampleKeys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "worker": ""
  }
 }
]
```
#### ------------ Add Traintuple ------------
Smart contract: `createTraintuple`

##### JSON Inputs:
```go
{
 "algoKey": string (required,len=64,hexadecimal),
 "inModels": [string] (omitempty,dive,len=64,hexadecimal),
 "dataManagerKey": string (required,len=64,hexadecimal),
 "dataSampleKeys": [string] (required,unique,gt=0,dive,len=64,hexadecimal),
 "computePlanID": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"inModels\":[],\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"computePlanID\":\"\",\"rank\":\"\",\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c"
}
```
#### ------------ Add Traintuple with inModel from previous traintuple ------------
Smart contract: `createTraintuple`

##### JSON Inputs:
```go
{
 "algoKey": string (required,len=64,hexadecimal),
 "inModels": [string] (omitempty,dive,len=64,hexadecimal),
 "dataManagerKey": string (required,len=64,hexadecimal),
 "dataSampleKeys": [string] (required,unique,gt=0,dive,len=64,hexadecimal),
 "computePlanID": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"inModels\":[\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"],\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"computePlanID\":\"\",\"rank\":\"\",\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d"
}
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d"
}
```
#### ------------ Query Traintuples of worker with todo status ------------
Smart contract: `queryFilter`

##### JSON Inputs:
```go
{
 "indexName": string (required),
 "attributes": string (required),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","{\"indexName\":\"traintuple~worker~status\",\"attributes\":\"SampleOrg, todo\"}"]}' -C myc
```
##### Command output:
```json
[
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "inModels": null,
  "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
  "log": "",
  "outModel": null,
  "permissions": {
   "process": {
    "authorizedIDs": [],
    "public": true
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": ""
 }
]
```
#### ------------ Log Start Training ------------
Smart contract: `logStartTrain`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "computePlanID": "",
 "creator": "SampleOrg",
 "dataset": {
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "inModels": null,
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
 "log": "",
 "outModel": null,
 "permissions": {
  "process": {
   "authorizedIDs": [],
   "public": true
  }
 },
 "rank": 0,
 "status": "doing",
 "tag": ""
}
```
#### ------------ Log Success Training ------------
Smart contract: `logSuccessTrain`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
 "log": string (lte=200),
 "outModel": (required){
   "hash": string (required,len=64,hexadecimal),
   "storageAddress": string (required),
 },
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\",\"log\":\"no error, ah ah ah\",\"outModel\":{\"hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"storageAddress\":\"https://substrabac/model/toto\"}}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "computePlanID": "",
 "creator": "SampleOrg",
 "dataset": {
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "inModels": null,
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
 "log": "no error, ah ah ah",
 "outModel": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto"
 },
 "permissions": {
  "process": {
   "authorizedIDs": [],
   "public": true
  }
 },
 "rank": 0,
 "status": "done",
 "tag": ""
}
```
#### ------------ Query Traintuple From key ------------
Smart contract: `queryTraintuple`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryTraintuple","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "computePlanID": "",
 "creator": "SampleOrg",
 "dataset": {
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "inModels": null,
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
 "log": "no error, ah ah ah",
 "outModel": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto"
 },
 "permissions": {
  "process": {
   "authorizedIDs": [],
   "public": true
  }
 },
 "rank": 0,
 "status": "done",
 "tag": ""
}
```
#### ------------ Add Non-Certified Testtuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "dataManagerKey": string (omitempty,len=64,hexadecimal),
 "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "objectiveKey": string (required,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
 "traintupleKey": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleKey\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8"
}
```
#### ------------ Add Certified Testtuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "dataManagerKey": string (omitempty,len=64,hexadecimal),
 "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "objectiveKey": string (required,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
 "traintupleKey": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleKey\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d"
}
```
#### ------------ Add Testtuple with not done traintuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "dataManagerKey": string (omitempty,len=64,hexadecimal),
 "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "objectiveKey": string (required,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
 "traintupleKey": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleKey\":\"ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "f5515eef9906b6af129355146a648c94f390e32d0fb45a8f54fddfe3329df716"
}
```
#### ------------ Query Testtuples of worker with todo status ------------
Smart contract: `queryFilter`

##### JSON Inputs:
```go
{
 "indexName": string (required),
 "attributes": string (required),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","{\"indexName\":\"testtuple~worker~status\",\"attributes\":\"SampleOrg, todo\"}"]}' -C myc
```
##### Command output:
```json
[
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": false,
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8",
  "log": "",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
  "traintupleType": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": true,
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
  "log": "",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
  "traintupleType": "traintuple"
 }
]
```
#### ------------ Log Start Testing ------------
Smart contract: `logStartTest`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","{\"key\":\"a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "certified": true,
 "computePlanID": "",
 "creator": "SampleOrg",
 "dataset": {
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0,
  "worker": "SampleOrg"
 },
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
 "log": "",
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "doing",
 "tag": "",
 "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
 "traintupleType": "traintuple"
}
```
#### ------------ Log Success Testing ------------
Smart contract: `logSuccessTest`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
 "log": string (lte=200),
 "perf": float32 (omitempty),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","{\"key\":\"a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d\",\"log\":\"no error, ah ah ah\",\"perf\":0.9}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "certified": true,
 "computePlanID": "",
 "creator": "SampleOrg",
 "dataset": {
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
 "log": "no error, ah ah ah",
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "done",
 "tag": "",
 "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
 "traintupleType": "traintuple"
}
```
#### ------------ Query Testtuple from its key ------------
Smart contract: `queryTesttuple`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryTesttuple","{\"key\":\"a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "certified": true,
 "computePlanID": "",
 "creator": "SampleOrg",
 "dataset": {
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
 "log": "no error, ah ah ah",
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "done",
 "tag": "",
 "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
 "traintupleType": "traintuple"
}
```
#### ------------ Query all Testtuples ------------
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryTesttuples"]}' -C myc
```
##### Command output:
```json
[
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": false,
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8",
  "log": "",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
  "traintupleType": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": true,
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
  "log": "no error, ah ah ah",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "done",
  "tag": "",
  "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
  "traintupleType": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": true,
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "f5515eef9906b6af129355146a648c94f390e32d0fb45a8f54fddfe3329df716",
  "log": "",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "waiting",
  "tag": "",
  "traintupleKey": "ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d",
  "traintupleType": "traintuple"
 }
]
```
#### ------------ Query details about a model ------------
Smart contract: `queryModelDetails`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryModelDetails","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "nonCertifiedTesttuples": [
  {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "certified": false,
   "computePlanID": "",
   "creator": "SampleOrg",
   "dataset": {
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0,
    "worker": "SampleOrg"
   },
   "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8",
   "log": "",
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "rank": 0,
   "status": "todo",
   "tag": "",
   "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
   "traintupleType": "traintuple"
  }
 ],
 "testtuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": true,
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
  "log": "no error, ah ah ah",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "done",
  "tag": "",
  "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
  "traintupleType": "traintuple"
 },
 "traintuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "computePlanID": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "inModels": null,
  "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
  "log": "no error, ah ah ah",
  "outModel": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto"
  },
  "permissions": {
   "process": {
    "authorizedIDs": [],
    "public": true
   }
  },
  "rank": 0,
  "status": "done",
  "tag": ""
 }
}
```
#### ------------ Query all models ------------
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryModels"]}' -C myc
```
##### Command output:
```json
[
 {
  "traintuple": {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "computePlanID": "",
   "creator": "SampleOrg",
   "dataset": {
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "worker": "SampleOrg"
   },
   "inModels": null,
   "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
   "log": "no error, ah ah ah",
   "outModel": {
    "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
    "storageAddress": "https://substrabac/model/toto"
   },
   "permissions": {
    "process": {
     "authorizedIDs": [],
     "public": true
    }
   },
   "rank": 0,
   "status": "done",
   "tag": ""
  }
 },
 {
  "traintuple": {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "computePlanID": "",
   "creator": "SampleOrg",
   "dataset": {
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "worker": "SampleOrg"
   },
   "inModels": [
    {
     "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
     "storageAddress": "https://substrabac/model/toto",
     "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c"
    }
   ],
   "key": "ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d",
   "log": "",
   "outModel": null,
   "permissions": {
    "process": {
     "authorizedIDs": [],
     "public": true
    }
   },
   "rank": 0,
   "status": "todo",
   "tag": ""
  }
 }
]
```
#### ------------ Query model permissions ------------
Smart contract: `queryModelPermissions`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryModelPermissions","{\"key\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "process": {
  "authorizedIDs": [],
  "public": true
 }
}
```
#### ------------ Query Dataset ------------
Smart contract: `queryDataset`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryDataset","{\"key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "description": {
  "hash": "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee",
  "storageAddress": "https://toto/dataManager/42234/description"
 },
 "key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
 "name": "liver slide",
 "objectiveKey": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
 "opener": {
  "hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "storageAddress": "https://toto/dataManager/42234/opener"
 },
 "owner": "SampleOrg",
 "permissions": {
  "process": {
   "authorizedIDs": [],
   "public": true
  }
 },
 "testDataSampleKeys": [
  "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
 ],
 "trainDataSampleKeys": [
  "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
 ],
 "type": "images"
}
```
#### ------------ Query nodes ------------
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryNodes"]}' -C myc
```
##### Command output:
```json
[
 {
  "id": "SampleOrg"
 }
]
```
#### ------------ Update Data Sample with new data manager ------------
Smart contract: `updateDataSample`

##### JSON Inputs:
```go
{
 "hashes": [string] (required,dive,len=64,hexadecimal),
 "dataManagerKeys": [string] (required,dive,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateDataSample","{\"hashes\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"dataManagerKeys\":[\"38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee\"]}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "{\"keys\": [\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"]}"
}
```
#### ------------ Query the new Dataset ------------
Smart contract: `queryDataset`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryDataset","{\"key\":\"38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "description": {
  "hash": "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee",
  "storageAddress": "https://toto/dataManager/42234/description"
 },
 "key": "38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee",
 "name": "liver slide",
 "objectiveKey": "",
 "opener": {
  "hash": "38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee",
  "storageAddress": "https://toto/dataManager/42234/opener"
 },
 "owner": "SampleOrg",
 "permissions": {
  "process": {
   "authorizedIDs": [],
   "public": true
  }
 },
 "testDataSampleKeys": [],
 "trainDataSampleKeys": [
  "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
 ],
 "type": "images"
}
```
#### ------------ Create a ComputePlan ------------
Smart contract: `createComputePlan`

##### JSON Inputs:
```go
{
 "tag": string (omitempty,lte=64),
 "traintuples": (omitempty) [{
   "dataManagerKey": string (required,len=64,hexadecimal),
   "dataSampleKeys": [string] (required,dive,len=64,hexadecimal),
   "algoKey": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "inModelsIDs": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
 }],
 "aggregatetuples": (omitempty) [{
   "algoKey": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "inModelsIDs": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "worker": string (required),
 }],
 "compositeTraintuples": (omitempty) [{
   "dataManagerKey": string (required,len=64,hexadecimal),
   "dataSampleKeys": [string] (required,dive,len=64,hexadecimal),
   "algoKey": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "inHeadModelID": string (required_with=InTrunkModelID,omitempty,len=64,hexadecimal),
   "inTrunkModelID": string (required_with=InHeadModelID,omitempty,len=64,hexadecimal),
   "OutTrunkModelPermissions": (required){
     "process": (required){
       "public": bool (required),
       "authorizedIDs": [string] (required),
     },
   },
   "tag": string (omitempty,lte=64),
 }],
 "testtuples": (omitempty) [{
   "dataManagerKey": string (omitempty,len=64,hexadecimal),
   "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
   "objectiveKey": string (required,len=64,hexadecimal),
   "tag": string (omitempty,lte=64),
   "traintupleID": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createComputePlan","{\"cleanModels\":false,\"tag\":\"a tag is simply a string\",\"traintuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"firstTraintupleID\",\"inModelsIDs\":null,\"tag\":\"\"},{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"secondTraintupleID\",\"inModelsIDs\":[\"firstTraintupleID\"],\"tag\":\"\"}],\"aggregatetuples\":null,\"compositeTraintuples\":null,\"testtuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleID\":\"secondTraintupleID\"}]}"]}' -C myc
```
##### Command output:
```json
{
 "IDToKey": {
  "firstTraintupleID": "a3cbf90bf2945cc7728f27894acf8e609fd65ef93340157f90706bd441d5f1b8",
  "secondTraintupleID": "61f210813e18f3a7af124ff32d6f5a46a276d5ae9fdd5d5c6f49a8e49c88366d"
 },
 "__metrics__": {
  "duration": 0
 },
 "aggregatetupleKeys": null,
 "cleanModels": false,
 "compositeTraintupleKeys": null,
 "computePlanID": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "doneCount": 0,
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtupleKeys": [
  "57c5c39799167b74ebe67974776b4149f5cc4ddde0e5811e75a19348191e2194"
 ],
 "traintupleKeys": [
  "a3cbf90bf2945cc7728f27894acf8e609fd65ef93340157f90706bd441d5f1b8",
  "61f210813e18f3a7af124ff32d6f5a46a276d5ae9fdd5d5c6f49a8e49c88366d"
 ],
 "tupleCount": 3
}
```
#### ------------ Update a ComputePlan ------------
Smart contract: `updateComputePlan`

##### JSON Inputs:
```go
{
 "computePlanID": string (required,required,len=64,hexadecimal),
 "traintuples": (omitempty) [{
   "dataManagerKey": string (required,len=64,hexadecimal),
   "dataSampleKeys": [string] (required,dive,len=64,hexadecimal),
   "algoKey": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "inModelsIDs": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
 }],
 "aggregatetuples": (omitempty) [{
   "algoKey": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "inModelsIDs": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "worker": string (required),
 }],
 "compositeTraintuples": (omitempty) [{
   "dataManagerKey": string (required,len=64,hexadecimal),
   "dataSampleKeys": [string] (required,dive,len=64,hexadecimal),
   "algoKey": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "inHeadModelID": string (required_with=InTrunkModelID,omitempty,len=64,hexadecimal),
   "inTrunkModelID": string (required_with=InHeadModelID,omitempty,len=64,hexadecimal),
   "OutTrunkModelPermissions": (required){
     "process": (required){
       "public": bool (required),
       "authorizedIDs": [string] (required),
     },
   },
   "tag": string (omitempty,lte=64),
 }],
 "testtuples": (omitempty) [{
   "dataManagerKey": string (omitempty,len=64,hexadecimal),
   "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
   "objectiveKey": string (required,len=64,hexadecimal),
   "tag": string (omitempty,lte=64),
   "traintupleID": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateComputePlan","{\"computePlanID\":\"7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17\",\"traintuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"thirdTraintupleID\",\"inModelsIDs\":[\"firstTraintupleID\",\"secondTraintupleID\"],\"tag\":\"\"}],\"aggregatetuples\":null,\"compositeTraintuples\":null,\"testtuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleID\":\"thirdTraintupleID\"}]}"]}' -C myc
```
##### Command output:
```json
{
 "IDToKey": {
  "thirdTraintupleID": "57a09687b8d26cc553eeac739fa6cc1c5e7f1a52a2ce2df277576eee7356282b"
 },
 "__metrics__": {
  "duration": 0
 },
 "aggregatetupleKeys": null,
 "cleanModels": false,
 "compositeTraintupleKeys": null,
 "computePlanID": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "doneCount": 0,
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtupleKeys": [
  "57c5c39799167b74ebe67974776b4149f5cc4ddde0e5811e75a19348191e2194",
  "2c78d12bd921a74f057ba5e535c9fda224e551e9562de863fd729cf347ded369"
 ],
 "traintupleKeys": [
  "a3cbf90bf2945cc7728f27894acf8e609fd65ef93340157f90706bd441d5f1b8",
  "61f210813e18f3a7af124ff32d6f5a46a276d5ae9fdd5d5c6f49a8e49c88366d",
  "57a09687b8d26cc553eeac739fa6cc1c5e7f1a52a2ce2df277576eee7356282b"
 ],
 "tupleCount": 5
}
```
#### ------------ Query an ObjectiveLeaderboard ------------
Smart contract: `queryObjectiveLeaderboard`

##### JSON Inputs:
```go
{
 "objectiveKey": string (omitempty,len=64,hexadecimal),
 "ascendingOrder": bool (required),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryObjectiveLeaderboard","{\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"ascendingOrder\":true}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "objective": {
  "description": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/description"
  },
  "key": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "name": "accuracy",
   "storageAddress": "https://toto/objective/222/metrics"
  },
  "name": "MSI classification",
  "owner": "SampleOrg",
  "permissions": {
   "process": {
    "authorizedIDs": [],
    "public": true
   }
  },
  "testDataset": {
   "dataManagerKey": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "dataSampleKeys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "worker": ""
  }
 },
 "testtuples": [
  {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "creator": "SampleOrg",
   "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
   "perf": 0.9,
   "tag": "",
   "traintupleKey": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c"
  }
 ]
}
```
#### ------------ Query Compute Plan(s) ------------
Smart contract: `queryComputePlan`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryComputePlan","{\"key\":\"7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17\"}"]}' -C myc
```
##### Command output:
```json
{
 "IDToKey": {},
 "__metrics__": {
  "duration": 0
 },
 "aggregatetupleKeys": null,
 "cleanModels": false,
 "compositeTraintupleKeys": null,
 "computePlanID": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "doneCount": 0,
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtupleKeys": [
  "57c5c39799167b74ebe67974776b4149f5cc4ddde0e5811e75a19348191e2194",
  "2c78d12bd921a74f057ba5e535c9fda224e551e9562de863fd729cf347ded369"
 ],
 "traintupleKeys": [
  "a3cbf90bf2945cc7728f27894acf8e609fd65ef93340157f90706bd441d5f1b8",
  "61f210813e18f3a7af124ff32d6f5a46a276d5ae9fdd5d5c6f49a8e49c88366d",
  "57a09687b8d26cc553eeac739fa6cc1c5e7f1a52a2ce2df277576eee7356282b"
 ],
 "tupleCount": 5
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryComputePlans"]}' -C myc
```
##### Command output:
```json
[
 {
  "IDToKey": {},
  "aggregatetupleKeys": null,
  "cleanModels": false,
  "compositeTraintupleKeys": null,
  "computePlanID": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
  "doneCount": 0,
  "status": "todo",
  "tag": "a tag is simply a string",
  "testtupleKeys": [
   "57c5c39799167b74ebe67974776b4149f5cc4ddde0e5811e75a19348191e2194",
   "2c78d12bd921a74f057ba5e535c9fda224e551e9562de863fd729cf347ded369"
  ],
  "traintupleKeys": [
   "a3cbf90bf2945cc7728f27894acf8e609fd65ef93340157f90706bd441d5f1b8",
   "61f210813e18f3a7af124ff32d6f5a46a276d5ae9fdd5d5c6f49a8e49c88366d",
   "57a09687b8d26cc553eeac739fa6cc1c5e7f1a52a2ce2df277576eee7356282b"
  ],
  "tupleCount": 5
 }
]
```
#### ------------ Cancel a ComputePlan ------------
Smart contract: `cancelComputePlan`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["cancelComputePlan","{\"key\":\"7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17\"}"]}' -C myc
```
##### Command output:
```json
{
 "IDToKey": {},
 "__metrics__": {
  "duration": 0
 },
 "aggregatetupleKeys": null,
 "cleanModels": false,
 "compositeTraintupleKeys": null,
 "computePlanID": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "doneCount": 0,
 "status": "canceled",
 "tag": "a tag is simply a string",
 "testtupleKeys": [
  "57c5c39799167b74ebe67974776b4149f5cc4ddde0e5811e75a19348191e2194",
  "2c78d12bd921a74f057ba5e535c9fda224e551e9562de863fd729cf347ded369"
 ],
 "traintupleKeys": [
  "a3cbf90bf2945cc7728f27894acf8e609fd65ef93340157f90706bd441d5f1b8",
  "61f210813e18f3a7af124ff32d6f5a46a276d5ae9fdd5d5c6f49a8e49c88366d",
  "57a09687b8d26cc553eeac739fa6cc1c5e7f1a52a2ce2df277576eee7356282b"
 ],
 "tupleCount": 5
}
```
