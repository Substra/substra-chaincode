# substra-chaincode
Chaincode for the Substra platform

> :warning: This project is under active development. Please, wait some times before using it...

## License

This project is developed under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](./LICENSE) file.
## Devmode

See [chaincode-docker-devmode](./chaincode-docker-devmode/README.rst)

## Documentation


Note for internal use only: See the [technical specifications](https://github.com/SubstraFoundation/substra-spec/blob/master/technical_spec_substra.md#smartcontract).

### Implemented smart contracts

- `createComputePlan`
- `createTesttuple`
- `createTraintuple`
- `logFailTest`
- `logFailTrain`
- `logStartTest`
- `logStartTrain`
- `logSuccessTest`
- `logSuccessTrain`
- `queryAlgo`
- `queryAlgos`
- `queryDataManager`
- `queryDataManagers`
- `queryDataset`
- `queryFilter`
- `queryModelDetails`
- `queryModels`
- `queryObjective`
- `queryObjectives`
- `queryTesttuple`
- `queryTesttuples`
- `queryTraintuple`
- `queryTraintuples`
- `registerAlgo`
- `registerDataManager`
- `registerDataSample`
- `registerObjective`
- `updateDataManager`
- `updateDataSample`
- `registerNode`
- `queryNodes`

### Examples

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
 "keys": [
  "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
 ]
}
```
#### ------------ Add Node ------------
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerNode"]}' -C myc
```
##### Command output:
```json
{
 "id": "SampleOrg"
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
 "objectiveKey": string (required,len=64,hexadecimal),
 "inModels": [string] (omitempty,dive,len=64,hexadecimal),
 "dataManagerKey": string (required,len=64,hexadecimal),
 "dataSampleKeys": [string] (required,gt=1,dive,len=64,hexadecimal),
 "computePlanID": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"inModels\":[],\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"computePlanID\":\"\",\"rank\":\"\",\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
}
```
#### ------------ Add Traintuple with inModel from previous traintuple ------------
Smart contract: `createTraintuple`

##### JSON Inputs:
```go
{
 "algoKey": string (required,len=64,hexadecimal),
 "objectiveKey": string (required,len=64,hexadecimal),
 "inModels": [string] (omitempty,dive,len=64,hexadecimal),
 "dataManagerKey": string (required,len=64,hexadecimal),
 "dataSampleKeys": [string] (required,gt=1,dive,len=64,hexadecimal),
 "computePlanID": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"inModels\":[\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\"],\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"computePlanID\":\"\",\"rank\":\"\",\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab"
}
```
##### Command output:
```json
{
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
   "perf": 0,
   "worker": "SampleOrg"
  },
  "inModels": null,
  "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
  "log": "",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
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
  "perf": 0,
  "worker": "SampleOrg"
 },
 "inModels": null,
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
 "log": "",
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
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
 "log": string (required,lte=200),
 "outModel": (required){
   "hash": string (required,len=64,hexadecimal),
   "storageAddress": string (required),
 },
 "perf": float32 (omitempty),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","{\"key\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\",\"log\":\"no error, ah ah ah\",\"outModel\":{\"hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"storageAddress\":\"https://substrabac/model/toto\"},\"perf\":0.9}"]}' -C myc
```
##### Command output:
```json
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
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "inModels": null,
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
 "log": "no error, ah ah ah",
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
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
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "inModels": null,
 "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
 "log": "no error, ah ah ah",
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
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
 "traintupleKey": string (required,len=64,hexadecimal),
 "dataManagerKey": string (omitempty,len=64,hexadecimal),
 "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"traintupleKey\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\",\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "c5f71e7a53c8a88af3e9b0311eaec68abd30718a388e8f8b45b0547ef2289dcd"
}
```
#### ------------ Add Certified Testtuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "traintupleKey": string (required,len=64,hexadecimal),
 "dataManagerKey": string (omitempty,len=64,hexadecimal),
 "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"traintupleKey\":\"9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3\",\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba"
}
```
#### ------------ Add Testtuple with not done traintuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "traintupleKey": string (required,len=64,hexadecimal),
 "dataManagerKey": string (omitempty,len=64,hexadecimal),
 "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"traintupleKey\":\"720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab\",\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "d009acea2d213bc7149ee15b0eb23217e7f06154b79c7046a73eb13a50c3f9dc"
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
  "certified": true,
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
  "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
  "log": "",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "status": "todo",
  "tag": ""
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": false,
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
  "key": "c5f71e7a53c8a88af3e9b0311eaec68abd30718a388e8f8b45b0547ef2289dcd",
  "log": "",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "status": "todo",
  "tag": ""
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","{\"key\":\"5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "certified": true,
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
 "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
 "log": "",
 "model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto",
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
 },
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "status": "doing",
 "tag": ""
}
```
#### ------------ Log Success Testing ------------
Smart contract: `logSuccessTest`

##### JSON Inputs:
```go
{
 "key": string (required,len=64,hexadecimal),
 "log": string (required,lte=200),
 "perf": float32 (omitempty),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","{\"key\":\"5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba\",\"log\":\"no error, ah ah ah\",\"perf\":0.9}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "certified": true,
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
 "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
 "log": "no error, ah ah ah",
 "model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto",
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
 },
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "status": "done",
 "tag": ""
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
peer chaincode query -n mycc -c '{"Args":["queryTesttuple","{\"key\":\"5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "certified": true,
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
 "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
 "log": "no error, ah ah ah",
 "model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto",
  "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
 },
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "status": "done",
 "tag": ""
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
  "key": "d009acea2d213bc7149ee15b0eb23217e7f06154b79c7046a73eb13a50c3f9dc",
  "log": "",
  "model": {
   "hash": "",
   "storageAddress": "",
   "traintupleKey": "720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "status": "waiting",
  "tag": ""
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": false,
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
  "key": "c5f71e7a53c8a88af3e9b0311eaec68abd30718a388e8f8b45b0547ef2289dcd",
  "log": "",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "status": "todo",
  "tag": ""
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": true,
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
  "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
  "log": "no error, ah ah ah",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "status": "done",
  "tag": ""
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
 "nonCertifiedTesttuples": [
  {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "certified": false,
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
   "key": "c5f71e7a53c8a88af3e9b0311eaec68abd30718a388e8f8b45b0547ef2289dcd",
   "log": "",
   "model": {
    "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
    "storageAddress": "https://substrabac/model/toto",
    "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
   },
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "status": "todo",
   "tag": ""
  }
 ],
 "testtuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "certified": true,
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
  "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
  "log": "no error, ah ah ah",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "status": "done",
  "tag": ""
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
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "inModels": null,
  "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
  "log": "no error, ah ah ah",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
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
  "testtuple": {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "certified": true,
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
   "key": "d009acea2d213bc7149ee15b0eb23217e7f06154b79c7046a73eb13a50c3f9dc",
   "log": "",
   "model": {
    "hash": "",
    "storageAddress": "",
    "traintupleKey": "720f778397fa07e24c2f314599725bf97727ded07ff65a51fa1a97b24d11ecab"
   },
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "status": "waiting",
   "tag": ""
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
    "perf": 0,
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
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
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
  "testtuple": {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "certified": true,
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
   "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
   "log": "no error, ah ah ah",
   "model": {
    "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
    "storageAddress": "https://substrabac/model/toto",
    "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
   },
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "status": "done",
   "tag": ""
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
    "perf": 0.9,
    "worker": "SampleOrg"
   },
   "inModels": null,
   "key": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3",
   "log": "no error, ah ah ah",
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
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
{
 "nodeIDs": [
  "SampleOrg"
 ]
}
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
 "algoKey": string (required,len=64,hexadecimal),
 "objectiveKey": string (required,len=64,hexadecimal),
 "traintuples": (required,gt=0) [{
   "dataManagerKey": string (required,len=64,hexadecimal),
   "dataSampleKeys": [string] (required,dive,len=64,hexadecimal),
   "id": string (required,lte=64),
   "inModelsIDs": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
 }],
 "testtuples": (omitempty) [{
   "dataManagerKey": string (omitempty,len=64,hexadecimal),
   "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
   "tag": string (omitempty,lte=64),
   "traintupleID": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createComputePlan","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"traintuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"id\":\"firstTraintupleID\",\"inModelsIDs\":null,\"tag\":\"\"},{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"id\":\"secondTraintupleID\",\"inModelsIDs\":[\"firstTraintupleID\"],\"tag\":\"\"}],\"testtuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"tag\":\"\",\"traintupleID\":\"secondTraintupleID\"}]}"]}' -C myc
```
##### Command output:
```json
{
 "computePlanID": "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
 "testtupleKeys": [
  "1dbd49d84e00ad6f339f416af0decfaf2db8f14412786de65b597e49a6820f96"
 ],
 "traintupleKeys": [
  "432fcffdf68892f5e4adeeed8bb618beaeaecf709f840671eca724a3e3109369",
  "d23f8cf290b902417ae698d68e2c6835483521d54fcbece31208517759b7c299"
 ]
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
   "key": "5ae68332a1e7182d9286692a892c7bf6f339d71d393ec6308e598c159d369aba",
   "model": {
    "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
    "storageAddress": "https://substrabac/model/toto",
    "traintupleKey": "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
   },
   "perf": 0.9,
   "tag": ""
  }
 ]
}
```
