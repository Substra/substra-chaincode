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
 "permissions": string (required,oneof=all),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataManager","{\"name\":\"liver slide\",\"openerHash\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"openerStorageAddress\":\"https://toto/dataManager/42234/opener\",\"type\":\"images\",\"descriptionHash\":\"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee\",\"descriptionStorageAddress\":\"https://toto/dataManager/42234/description\",\"objectiveKey\":\"\",\"permissions\":\"all\"}"]}' -C myc
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
 "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "permissions": "all",
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
   }
  "permissions": string (required,oneof=all),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerObjective","{\"name\":\"MSI classification\",\"descriptionHash\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"descriptionStorageAddress\":\"https://toto/objective/222/description\",\"metricsName\":\"accuracy\",\"metricsHash\":\"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"metricsStorageAddress\":\"https://toto/objective/222/metrics\",\"testDataset\":{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"]},\"permissions\":\"all\"}"]}' -C myc
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
 "permissions": string (required,oneof=all),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo","{\"name\":\"hog + svm\",\"hash\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"storageAddress\":\"https://toto/algo/222/algo\",\"descriptionHash\":\"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca\",\"descriptionStorageAddress\":\"https://toto/algo/222/description\",\"permissions\":\"all\"}"]}' -C myc
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
  "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "permissions": "all",
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
  "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 {
  "dataManagerKeys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 {
  "dataManagerKeys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 {
  "dataManagerKeys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
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
  "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "permissions": "all",
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
 "flTask": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"inModels\":[],\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"flTask\":\"\",\"rank\":\"\",\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
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
 "flTask": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"inModels\":[\"4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6\"],\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"flTask\":\"\",\"rank\":\"\",\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "05a8b7fa61a4f77f126f569b4ea108db20e35904f704dc233b88962375c46fc0"
}
```
##### Command output:
```json
{
 "key": "05a8b7fa61a4f77f126f569b4ea108db20e35904f704dc233b88962375c46fc0"
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
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","{\"indexName\":\"traintuple~worker~status\",\"attributes\":\"ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776, todo\"}"]}' -C myc
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
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "fltask": "",
  "inModels": null,
  "key": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6",
  "log": "",
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "outModel": null,
  "permissions": "all",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain","{\"key\":\"4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "dataset": {
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0,
  "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 "fltask": "",
 "inModels": null,
 "key": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6",
 "log": "",
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "outModel": null,
 "permissions": "all",
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
   }
  "perf": float32 (omitempty),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","{\"key\":\"4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6\",\"log\":\"no error, ah ah ah\",\"outModel\":{\"hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"storageAddress\":\"https://substrabac/model/toto\"},\"perf\":0.9}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "dataset": {
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 "fltask": "",
 "inModels": null,
 "key": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6",
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
 "permissions": "all",
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
peer chaincode invoke -n mycc -c '{"Args":["queryTraintuple","{\"key\":\"4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storageAddress": "https://toto/algo/222/algo"
 },
 "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "dataset": {
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 "fltask": "",
 "inModels": null,
 "key": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6",
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
 "permissions": "all",
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
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"traintupleKey\":\"4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6\",\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "e24b10100ae235aa0955ae42016d4dde3c48d1e955d0f9672e6e692d5f5e6c6e"
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
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"traintupleKey\":\"4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6\",\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339"
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
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"traintupleKey\":\"05a8b7fa61a4f77f126f569b4ea108db20e35904f704dc233b88962375c46fc0\",\"dataManagerKey\":\"\",\"dataSampleKeys\":null,\"tag\":\"\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "412d55375a3f19c8b11a4345a95a32c915f1346a9b65294ddf87908857ddbeab"
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
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","{\"indexName\":\"testtuple~worker~status\",\"attributes\":\"ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776, todo\"}"]}' -C myc
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
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339",
  "log": "",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "permissions": "all",
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
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "key": "e24b10100ae235aa0955ae42016d4dde3c48d1e955d0f9672e6e692d5f5e6c6e",
  "log": "",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "permissions": "all",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","{\"key\":\"96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339\"}"]}' -C myc
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
 "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "dataset": {
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0,
  "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339",
 "log": "",
 "model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto",
  "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
 },
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "permissions": "all",
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","{\"key\":\"96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339\",\"log\":\"no error, ah ah ah\",\"perf\":0.9}"]}' -C myc
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
 "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "dataset": {
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339",
 "log": "no error, ah ah ah",
 "model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto",
  "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
 },
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "permissions": "all",
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
peer chaincode query -n mycc -c '{"Args":["queryTesttuple","{\"key\":\"96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339\"}"]}' -C myc
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
 "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "dataset": {
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
 },
 "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339",
 "log": "no error, ah ah ah",
 "model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storageAddress": "https://substrabac/model/toto",
  "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
 },
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storageAddress": "https://toto/objective/222/metrics"
  }
 },
 "permissions": "all",
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
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "key": "412d55375a3f19c8b11a4345a95a32c915f1346a9b65294ddf87908857ddbeab",
  "log": "",
  "model": {
   "hash": "",
   "storageAddress": "",
   "traintupleKey": "05a8b7fa61a4f77f126f569b4ea108db20e35904f704dc233b88962375c46fc0"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "permissions": "all",
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
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "key": "e24b10100ae235aa0955ae42016d4dde3c48d1e955d0f9672e6e692d5f5e6c6e",
  "log": "",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "permissions": "all",
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
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339",
  "log": "no error, ah ah ah",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "permissions": "all",
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
peer chaincode query -n mycc -c '{"Args":["queryModelDetails","{\"key\":\"4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6\"}"]}' -C myc
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
   "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
   "dataset": {
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0,
    "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
   },
   "key": "e24b10100ae235aa0955ae42016d4dde3c48d1e955d0f9672e6e692d5f5e6c6e",
   "log": "",
   "model": {
    "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
    "storageAddress": "https://substrabac/model/toto",
    "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
   },
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "permissions": "all",
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
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339",
  "log": "no error, ah ah ah",
  "model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storageAddress": "https://substrabac/model/toto",
   "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
  },
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storageAddress": "https://toto/objective/222/metrics"
   }
  },
  "permissions": "all",
  "status": "done",
  "tag": ""
 },
 "traintuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storageAddress": "https://toto/algo/222/algo"
  },
  "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
  },
  "fltask": "",
  "inModels": null,
  "key": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6",
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
  "permissions": "all",
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
   "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
   "dataset": {
    "keys": [
     "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0,
    "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
   },
   "key": "412d55375a3f19c8b11a4345a95a32c915f1346a9b65294ddf87908857ddbeab",
   "log": "",
   "model": {
    "hash": "",
    "storageAddress": "",
    "traintupleKey": "05a8b7fa61a4f77f126f569b4ea108db20e35904f704dc233b88962375c46fc0"
   },
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "permissions": "all",
   "status": "waiting",
   "tag": ""
  },
  "traintuple": {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
   "dataset": {
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0,
    "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
   },
   "fltask": "",
   "inModels": [
    {
     "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
     "storageAddress": "https://substrabac/model/toto",
     "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
    }
   ],
   "key": "05a8b7fa61a4f77f126f569b4ea108db20e35904f704dc233b88962375c46fc0",
   "log": "",
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "outModel": null,
   "permissions": "all",
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
   "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
   "dataset": {
    "keys": [
     "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0.9,
    "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
   },
   "key": "96d9d7922847ef768db438d5cbfe1b09832d935e72133d93628d52dc81298339",
   "log": "no error, ah ah ah",
   "model": {
    "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
    "storageAddress": "https://substrabac/model/toto",
    "traintupleKey": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6"
   },
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storageAddress": "https://toto/objective/222/metrics"
    }
   },
   "permissions": "all",
   "status": "done",
   "tag": ""
  },
  "traintuple": {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storageAddress": "https://toto/algo/222/algo"
   },
   "creator": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
   "dataset": {
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "openerHash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0.9,
    "worker": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776"
   },
   "fltask": "",
   "inModels": null,
   "key": "4ef2002b70fdeb9d91679a98eb74a74edb7db852cd064209807f79d5dfd22cd6",
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
   "permissions": "all",
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
 "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "permissions": "all",
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
 "owner": "ec722969cb64200ab6638f68ac538e40abab5b19a6485661042a1061c4612776",
 "permissions": "all",
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
 "traintuples": (required,gt=0) [
   {
    "dataManagerKey": string (required,len=64,hexadecimal),
    "dataSampleKeys": [string] (required,dive,len=64,hexadecimal),
    "id": string (required,lte=64),
    "inModelsIDs": [string] (omitempty,dive,lte=64),
    "tag": string (omitempty,lte=64),
   }
 ],
 "testtuples": (omitempty) [
   {
    "dataManagerKey": string (omitempty,len=64,hexadecimal),
    "dataSampleKeys": [string] (omitempty,dive,len=64,hexadecimal),
    "tag": string (omitempty,lte=64),
    "traintupleID": string (required,lte=64),
   }
 ],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createComputePlan","{\"algoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"objectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"traintuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"id\":\"firstTraintupleID\",\"inModelsIDs\":null,\"tag\":\"\"},{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"id\":\"secondTraintupleID\",\"inModelsIDs\":[\"firstTraintupleID\"],\"tag\":\"\"}],\"testtuples\":[{\"dataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"dataSampleKeys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"tag\":\"\",\"traintupleID\":\"secondTraintupleID\"}]}"]}' -C myc
```
##### Command output:
```json
{
 "fltask": "60ad8fede3a6139b366008482ab9f29233b09f521d09ac3ad50cfcb55a8fadc9",
 "testtupleKeys": [
  "c5805fea88648f6abcddf5e6862133e4f8efbc67561b1380ece06eeec1623bd4"
 ],
 "traintupleKeys": [
  "60ad8fede3a6139b366008482ab9f29233b09f521d09ac3ad50cfcb55a8fadc9",
  "e95e6f91581c458e8c742535b86c872405f8c945da796ce91a311d6dd3e05d4c"
 ]
}
```
