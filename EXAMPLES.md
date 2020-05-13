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
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
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
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"inModels\":[\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\"],\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"computePlanID\":\"\",\"rank\":\"\",\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab"
}
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab"
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
  "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain","{\"key\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\"}"]}' -C myc
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
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","{\"key\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\",\"log\":\"no error, ah ah ah\",\"outModel\":{\"hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"storageAddress\":\"https://substrabac/model/toto\"}}"]}' -C myc
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
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode invoke -n mycc -c '{"Args":["queryTraintuple","{\"key\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\"}"]}' -C myc
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
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleKey\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "3b807eb0bcd6b0798dc8f6eb415d2e58fb4d3515d2b00d4b888be1ca8145b7d8"
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
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleKey\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e"
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
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"traintupleKey\":\"720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab\"}"]}' -C myc
```
##### Command output:
```json
{
 "__metrics__": {
  "duration": 0
 },
 "key": "76d4a5908359eb6ba9c8fc89254c4b08e23aa20471ea3ddeee9a2835825dbd72"
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
  "key": "3b807eb0bcd6b0798dc8f6eb415d2e58fb4d3515d2b00d4b888be1ca8145b7d8",
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
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
  "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e",
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
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","{\"key\":\"4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e\"}"]}' -C myc
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
 "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e",
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
 "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","{\"key\":\"4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e\",\"log\":\"no error, ah ah ah\",\"perf\":0.9}"]}' -C myc
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
 "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e",
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
 "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode query -n mycc -c '{"Args":["queryTesttuple","{\"key\":\"4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e\"}"]}' -C myc
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
 "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e",
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
 "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
  "key": "76d4a5908359eb6ba9c8fc89254c4b08e23aa20471ea3ddeee9a2835825dbd72",
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
  "traintupleKey": "720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab",
  "traintupleType": "traintuple"
 },
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
  "key": "3b807eb0bcd6b0798dc8f6eb415d2e58fb4d3515d2b00d4b888be1ca8145b7d8",
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
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
  "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e",
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
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
peer chaincode query -n mycc -c '{"Args":["queryModelDetails","{\"key\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\"}"]}' -C myc
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
   "key": "3b807eb0bcd6b0798dc8f6eb415d2e58fb4d3515d2b00d4b888be1ca8145b7d8",
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
   "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
  "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e",
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
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
  "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
   "inModels": [
    {
     "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
     "storageAddress": "https://substrabac/model/toto",
     "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
    }
   ],
   "key": "720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab",
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
   "inModels": null,
   "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
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
  "firstTraintupleID": "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
  "secondTraintupleID": "d23f8cf290b902417ae698d68e2c6835483521d54fcbece31208517759b7c299"
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
  "18285ef38b518c2ac73e9bffa3523c55e2ea6c968a7d7db89c9e8ab0a14c2562"
 ],
 "traintupleKeys": [
  "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
  "d23f8cf290b902417ae698d68e2c6835483521d54fcbece31208517759b7c299"
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
  "thirdTraintupleID": "c163663889566a4c51fad3765107774891f8ed456dcc859a9a1a7fcd17b11386"
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
  "18285ef38b518c2ac73e9bffa3523c55e2ea6c968a7d7db89c9e8ab0a14c2562",
  "a350eb73efd73b9797ba3a3d2d10f36145274176a3805fbd9f598540192f57f3"
 ],
 "traintupleKeys": [
  "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
  "d23f8cf290b902417ae698d68e2c6835483521d54fcbece31208517759b7c299",
  "c163663889566a4c51fad3765107774891f8ed456dcc859a9a1a7fcd17b11386"
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
   "key": "4d49bf9147bf391f9610d830aae6630290e128dacd7c3540e82178a0e002951e",
   "perf": 0.9,
   "tag": "",
   "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
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
  "18285ef38b518c2ac73e9bffa3523c55e2ea6c968a7d7db89c9e8ab0a14c2562",
  "a350eb73efd73b9797ba3a3d2d10f36145274176a3805fbd9f598540192f57f3"
 ],
 "traintupleKeys": [
  "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
  "d23f8cf290b902417ae698d68e2c6835483521d54fcbece31208517759b7c299",
  "c163663889566a4c51fad3765107774891f8ed456dcc859a9a1a7fcd17b11386"
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
   "18285ef38b518c2ac73e9bffa3523c55e2ea6c968a7d7db89c9e8ab0a14c2562",
   "a350eb73efd73b9797ba3a3d2d10f36145274176a3805fbd9f598540192f57f3"
  ],
  "traintupleKeys": [
   "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
   "d23f8cf290b902417ae698d68e2c6835483521d54fcbece31208517759b7c299",
   "c163663889566a4c51fad3765107774891f8ed456dcc859a9a1a7fcd17b11386"
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
  "18285ef38b518c2ac73e9bffa3523c55e2ea6c968a7d7db89c9e8ab0a14c2562",
  "a350eb73efd73b9797ba3a3d2d10f36145274176a3805fbd9f598540192f57f3"
 ],
 "traintupleKeys": [
  "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
  "d23f8cf290b902417ae698d68e2c6835483521d54fcbece31208517759b7c299",
  "c163663889566a4c51fad3765107774891f8ed456dcc859a9a1a7fcd17b11386"
 ],
 "tupleCount": 5
}
```
