### Examples

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
#### ------------ Add DataManager ------------
Smart contract: `registerDataManager`

##### JSON Inputs:
```go
{
 "key": string (required,len=36),
 "name": string (required,gte=1,lte=100),
 "opener_hash": string (required,len=64,hexadecimal),
 "opener_storage_address": string (required,url),
 "type": string (required,gte=1,lte=30),
 "description_hash": string (required,len=64,hexadecimal),
 "description_storage_address": string (required,url),
 "objective_key": string (omitempty,len=36),
 "permissions": (required){
   "process": (required){
     "public": bool (required),
     "authorized_ids": [string] (required),
   },
 },
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataManager","{\"key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"name\":\"liver slide\",\"opener_hash\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"opener_storage_address\":\"https://toto/dataManager/42234/opener\",\"type\":\"images\",\"description_hash\":\"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee\",\"description_storage_address\":\"https://toto/dataManager/42234/description\",\"objective_key\":\"\",\"permissions\":{\"process\":{\"public\":true,\"authorized_ids\":[]}},\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042"
}
```
#### ------------ Query DataManager From key ------------
Smart contract: `queryDataManager`

##### JSON Inputs:
```go
{
 "key": string (required,len=36),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryDataManager","{\"key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\"}"]}' -C myc
```
##### Command output:
```json
{
 "description": {
  "hash": "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee",
  "storage_address": "https://toto/dataManager/42234/description"
 },
 "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
 "metadata": {},
 "name": "liver slide",
 "objective_key": "",
 "opener": {
  "hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "storage_address": "https://toto/dataManager/42234/opener"
 },
 "owner": "SampleOrg",
 "permissions": {
  "process": {
   "authorized_ids": [],
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
 "keys": [string] (required,dive,len=36),
 "data_manager_keys": [string] (omitempty,dive,len=36),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"keys\":[\"bb1bb7c3-1f62-244c-0f3a-761cc1688042\",\"bb2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"data_manager_keys\":[\"da1bb7c3-1f62-244c-0f3a-761cc1688042\"],\"testOnly\":\"true\"}"]}' -C myc
```
##### Command output:
```json
{
 "keys": [
  "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
  "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
 ]
}
```
#### ------------ Add Objective ------------
Smart contract: `registerObjective`

##### JSON Inputs:
```go
{
 "key": string (required,len=36),
 "name": string (required,gte=1,lte=100),
 "description_hash": string (required,len=64,hexadecimal),
 "description_storage_address": string (required,url),
 "metrics_name": string (required,gte=1,lte=100),
 "metrics_hash": string (required,len=64,hexadecimal),
 "metrics_storage_address": string (required,url),
 "test_dataset": (omitempty){
   "data_manager_key": string (omitempty,len=36),
   "data_sample_keys": [string] (omitempty,dive,len=36),
 },
 "permissions": (required){
   "process": (required){
     "public": bool (required),
     "authorized_ids": [string] (required),
   },
 },
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerObjective","{\"key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"name\":\"MSI classification\",\"description_hash\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"description_storage_address\":\"https://toto/objective/222/description\",\"metrics_name\":\"accuracy\",\"metrics_hash\":\"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"metrics_storage_address\":\"https://toto/objective/222/metrics\",\"test_dataset\":{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"bb1bb7c3-1f62-244c-0f3a-761cc1688042\",\"bb2bb7c3-1f62-244c-0f3a-761cc1688042\"]},\"permissions\":{\"process\":{\"public\":true,\"authorized_ids\":[]}},\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "5c1d9cd1-c2c1-082d-de09-21b56d11030c"
}
```
#### ------------ Add Algo ------------
Smart contract: `registerAlgo`

##### JSON Inputs:
```go
{
 "key": string (required,len=36),
 "name": string (required,gte=1,lte=100),
 "hash": string (required,len=64,hexadecimal),
 "storage_address": string (required,url),
 "description_hash": string (required,len=64,hexadecimal),
 "description_storage_address": string (required,url),
 "permissions": (required){
   "process": (required){
     "public": bool (required),
     "authorized_ids": [string] (required),
   },
 },
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo","{\"key\":\"fd1bb7c3-1f62-244c-0f3a-761cc1688042\",\"name\":\"hog + svm\",\"hash\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"storage_address\":\"https://toto/algo/222/algo\",\"description_hash\":\"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca\",\"description_storage_address\":\"https://toto/algo/222/description\",\"permissions\":{\"process\":{\"public\":true,\"authorized_ids\":[]}},\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042"
}
```
#### ------------ Add Train DataSample ------------
Smart contract: `registerDataSample`

##### JSON Inputs:
```go
{
 "keys": [string] (required,dive,len=36),
 "data_manager_keys": [string] (omitempty,dive,len=36),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"keys\":[\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\",\"aa2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"data_manager_keys\":[\"da1bb7c3-1f62-244c-0f3a-761cc1688042\"],\"testOnly\":\"false\"}"]}' -C myc
```
##### Command output:
```json
{
 "keys": [
  "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
  "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
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
   "storage_address": "https://toto/dataManager/42234/description"
  },
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "metadata": {},
  "name": "liver slide",
  "objective_key": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
  "opener": {
   "hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "storage_address": "https://toto/dataManager/42234/opener"
  },
  "owner": "SampleOrg",
  "permissions": {
   "process": {
    "authorized_ids": [],
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
  "data_manager_keys": [
   "da1bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "key": "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "key": "aa2bb7c3-1f62-244c-0f3a-761cc1688042",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "key": "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "key": "bb2bb7c3-1f62-244c-0f3a-761cc1688042",
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
   "storage_address": "https://toto/objective/222/description"
  },
  "key": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
  "metadata": {},
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "name": "accuracy",
   "storage_address": "https://toto/objective/222/metrics"
  },
  "name": "MSI classification",
  "owner": "SampleOrg",
  "permissions": {
   "process": {
    "authorized_ids": [],
    "public": true
   }
  },
  "test_dataset": {
   "data_manager_key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "data_sample_keys": [
    "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
    "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "metadata": {},
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
 "algo_key": string (required,len=36),
 "in_models": [string] (omitempty,dive,len=64,hexadecimal),
 "data_manager_key": string (required,len=36),
 "data_sample_keys": [string] (required,unique,gt=0,dive,len=36),
 "compute_plan_id": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algo_key\":\"fd1bb7c3-1f62-244c-0f3a-761cc1688042\",\"in_models\":[],\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\",\"aa2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"compute_plan_id\":\"\",\"rank\":\"\",\"tag\":\"\",\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75"
}
```
#### ------------ Add Traintuple with inModel from previous traintuple ------------
Smart contract: `createTraintuple`

##### JSON Inputs:
```go
{
 "algo_key": string (required,len=36),
 "in_models": [string] (omitempty,dive,len=64,hexadecimal),
 "data_manager_key": string (required,len=36),
 "data_sample_keys": [string] (required,unique,gt=0,dive,len=36),
 "compute_plan_id": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algo_key\":\"fd1bb7c3-1f62-244c-0f3a-761cc1688042\",\"in_models\":[\"b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75\"],\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\",\"aa2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"compute_plan_id\":\"\",\"rank\":\"\",\"tag\":\"\",\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "860a31e0d50769bf11ec4424123febb8c2f813312c6104370a7e6cdf5644b53e"
}
```
##### Command output:
```json
{
 "key": "860a31e0d50769bf11ec4424123febb8c2f813312c6104370a7e6cdf5644b53e"
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
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
    "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "metadata": {},
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "in_models": null,
  "key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
  "log": "",
  "metadata": {},
  "out_model": null,
  "permissions": {
   "process": {
    "authorized_ids": [],
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain","{\"key\":\"b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
   "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
 "log": "",
 "metadata": {},
 "out_model": null,
 "permissions": {
  "process": {
   "authorized_ids": [],
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
 "out_model": (required){
   "hash": string (required,len=64,hexadecimal),
   "storage_address": string (required),
 },
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","{\"key\":\"b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75\",\"log\":\"no error, ah ah ah\",\"out_model\":{\"hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"storage_address\":\"https://substrabac/model/toto\"}}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
   "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
 "log": "no error, ah ah ah",
 "metadata": {},
 "out_model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storage_address": "https://substrabac/model/toto"
 },
 "permissions": {
  "process": {
   "authorized_ids": [],
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
peer chaincode invoke -n mycc -c '{"Args":["queryTraintuple","{\"key\":\"b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
   "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
 "log": "no error, ah ah ah",
 "metadata": {},
 "out_model": {
  "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
  "storage_address": "https://substrabac/model/toto"
 },
 "permissions": {
  "process": {
   "authorized_ids": [],
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
 "data_manager_key": string (omitempty,len=36),
 "data_sample_keys": [string] (omitempty,dive,len=36),
 "objective_key": string (required,len=36),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\",\"aa2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "449f44c60493e649bb386acc7cd49b2e87368cb258870287c284032b8b3525a6"
}
```
#### ------------ Add Certified Testtuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "data_manager_key": string (omitempty,len=36),
 "data_sample_keys": [string] (omitempty,dive,len=36),
 "objective_key": string (required,len=36),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"\",\"data_sample_keys\":null,\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883"
}
```
#### ------------ Add Testtuple with not done traintuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "data_manager_key": string (omitempty,len=36),
 "data_sample_keys": [string] (omitempty,dive,len=36),
 "objective_key": string (required,len=36),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"\",\"data_sample_keys\":null,\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"860a31e0d50769bf11ec4424123febb8c2f813312c6104370a7e6cdf5644b53e\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "f1de6499f5373178b8b8f914eaccdb210f4e9d43ab645aac16dfa86ce6fe9fc0"
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
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": false,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
    "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "449f44c60493e649bb386acc7cd49b2e87368cb258870287c284032b8b3525a6",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
  "traintuple_type": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
    "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
  "traintuple_type": "traintuple"
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","{\"key\":\"f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "certified": true,
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
   "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0,
  "worker": "SampleOrg"
 },
 "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883",
 "log": "",
 "metadata": {},
 "objective": {
  "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "doing",
 "tag": "",
 "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
 "traintuple_type": "traintuple"
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","{\"key\":\"f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883\",\"log\":\"no error, ah ah ah\",\"perf\":0.9}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "certified": true,
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
   "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883",
 "log": "no error, ah ah ah",
 "metadata": {},
 "objective": {
  "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "done",
 "tag": "",
 "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
 "traintuple_type": "traintuple"
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
peer chaincode query -n mycc -c '{"Args":["queryTesttuple","{\"key\":\"f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "certified": true,
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
   "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883",
 "log": "no error, ah ah ah",
 "metadata": {},
 "objective": {
  "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "done",
 "tag": "",
 "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
 "traintuple_type": "traintuple"
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
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
    "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "f1de6499f5373178b8b8f914eaccdb210f4e9d43ab645aac16dfa86ce6fe9fc0",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "waiting",
  "tag": "",
  "traintuple_key": "860a31e0d50769bf11ec4424123febb8c2f813312c6104370a7e6cdf5644b53e",
  "traintuple_type": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": false,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
    "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "449f44c60493e649bb386acc7cd49b2e87368cb258870287c284032b8b3525a6",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
  "traintuple_type": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
    "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883",
  "log": "no error, ah ah ah",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "done",
  "tag": "",
  "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
  "traintuple_type": "traintuple"
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
peer chaincode query -n mycc -c '{"Args":["queryModelDetails","{\"key\":\"b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75\"}"]}' -C myc
```
##### Command output:
```json
{
 "non_certified_testtuples": [
  {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
    "name": "hog + svm",
    "storage_address": "https://toto/algo/222/algo"
   },
   "certified": false,
   "compute_plan_id": "",
   "creator": "SampleOrg",
   "dataset": {
    "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
    "keys": [
     "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
     "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
    ],
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0,
    "worker": "SampleOrg"
   },
   "key": "449f44c60493e649bb386acc7cd49b2e87368cb258870287c284032b8b3525a6",
   "log": "",
   "metadata": {},
   "objective": {
    "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storage_address": "https://toto/objective/222/metrics"
    }
   },
   "rank": 0,
   "status": "todo",
   "tag": "",
   "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
   "traintuple_type": "traintuple"
  }
 ],
 "testtuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
    "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883",
  "log": "no error, ah ah ah",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "done",
  "tag": "",
  "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
  "traintuple_type": "traintuple"
 },
 "traintuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
    "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "metadata": {},
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "in_models": null,
  "key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
  "log": "no error, ah ah ah",
  "metadata": {},
  "out_model": {
   "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
   "storage_address": "https://substrabac/model/toto"
  },
  "permissions": {
   "process": {
    "authorized_ids": [],
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
    "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
    "name": "hog + svm",
    "storage_address": "https://toto/algo/222/algo"
   },
   "compute_plan_id": "",
   "creator": "SampleOrg",
   "dataset": {
    "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
    "keys": [
     "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
     "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
    ],
    "metadata": {},
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "worker": "SampleOrg"
   },
   "in_models": [
    {
     "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
     "storage_address": "https://substrabac/model/toto",
     "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75"
    }
   ],
   "key": "860a31e0d50769bf11ec4424123febb8c2f813312c6104370a7e6cdf5644b53e",
   "log": "",
   "metadata": {},
   "out_model": null,
   "permissions": {
    "process": {
     "authorized_ids": [],
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
    "key": "fd1bb7c3-1f62-244c-0f3a-761cc1688042",
    "name": "hog + svm",
    "storage_address": "https://toto/algo/222/algo"
   },
   "compute_plan_id": "",
   "creator": "SampleOrg",
   "dataset": {
    "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
    "keys": [
     "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
     "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
    ],
    "metadata": {},
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "worker": "SampleOrg"
   },
   "in_models": null,
   "key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75",
   "log": "no error, ah ah ah",
   "metadata": {},
   "out_model": {
    "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
    "storage_address": "https://substrabac/model/toto"
   },
   "permissions": {
    "process": {
     "authorized_ids": [],
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
 "process": {
  "authorized_ids": [],
  "public": true
 }
}
```
#### ------------ Query Dataset ------------
Smart contract: `queryDataset`

##### JSON Inputs:
```go
{
 "key": string (required,len=36),
}
```
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryDataset","{\"key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\"}"]}' -C myc
```
##### Command output:
```json
{
 "description": {
  "hash": "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee",
  "storage_address": "https://toto/dataManager/42234/description"
 },
 "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
 "metadata": {},
 "name": "liver slide",
 "objective_key": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
 "opener": {
  "hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "storage_address": "https://toto/dataManager/42234/opener"
 },
 "owner": "SampleOrg",
 "permissions": {
  "process": {
   "authorized_ids": [],
   "public": true
  }
 },
 "test_data_sample_keys": [
  "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
  "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
 ],
 "train_data_sample_keys": [
  "aa1bb7c3-1f62-244c-0f3a-761cc1688042",
  "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
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
 "keys": [string] (required,dive,len=36),
 "data_manager_keys": [string] (required,dive,len=36),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateDataSample","{\"keys\":[\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\"],\"data_manager_keys\":[\"38a320b2-a67c-8003-cc74-8d6666534f2b\"]}"]}' -C myc
```
##### Command output:
```json
{
 "key": "{\"keys\": [\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\"]}"
}
```
#### ------------ Query the new Dataset ------------
Smart contract: `queryDataset`

##### JSON Inputs:
```go
{
 "key": string (required,len=36),
}
```
##### Command peer example:
```bash
peer chaincode query -n mycc -c '{"Args":["queryDataset","{\"key\":\"38a320b2-a67c-8003-cc74-8d6666534f2b\"}"]}' -C myc
```
##### Command output:
```json
{
 "description": {
  "hash": "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee",
  "storage_address": "https://toto/dataManager/42234/description"
 },
 "key": "38a320b2-a67c-8003-cc74-8d6666534f2b",
 "metadata": {},
 "name": "liver slide",
 "objective_key": "",
 "opener": {
  "hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "storage_address": "https://toto/dataManager/42234/opener"
 },
 "owner": "SampleOrg",
 "permissions": {
  "process": {
   "authorized_ids": [],
   "public": true
  }
 },
 "test_data_sample_keys": [],
 "train_data_sample_keys": [
  "aa1bb7c3-1f62-244c-0f3a-761cc1688042"
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
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuples": (omitempty) [{
   "data_manager_key": string (required,len=36),
   "data_sample_keys": [string] (required,dive,len=36),
   "algo_key": string (required,len=36),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 }],
 "aggregatetuples": (omitempty) [{
   "algo_key": string (required,len=36),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "worker": string (required),
 }],
 "composite_traintuples": (omitempty) [{
   "data_manager_key": string (required,len=36),
   "data_sample_keys": [string] (required,dive,len=36),
   "algo_key": string (required,len=36),
   "id": string (required,lte=64),
   "in_head_model_id": string (required_with=InTrunkModelID,omitempty,len=64,hexadecimal),
   "in_trunk_model_id": string (required_with=InHeadModelID,omitempty,len=64,hexadecimal),
   "out_trunk_model_permissions": (required){
     "process": (required){
       "public": bool (required),
       "authorized_ids": [string] (required),
     },
   },
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 }],
 "testtuples": (omitempty) [{
   "data_manager_key": string (omitempty,len=36),
   "data_sample_keys": [string] (omitempty,dive,len=36),
   "objective_key": string (required,len=36),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "traintuple_id": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createComputePlan","{\"clean_models\":false,\"tag\":\"a tag is simply a string\",\"metadata\":null,\"traintuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\"],\"algo_key\":\"fd1bb7c3-1f62-244c-0f3a-761cc1688042\",\"id\":\"firstTraintupleID\",\"in_models_ids\":null,\"tag\":\"\",\"metadata\":null},{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"algo_key\":\"fd1bb7c3-1f62-244c-0f3a-761cc1688042\",\"id\":\"secondTraintupleID\",\"in_models_ids\":[\"firstTraintupleID\"],\"tag\":\"\",\"metadata\":null}],\"aggregatetuples\":null,\"composite_traintuples\":null,\"testtuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"bb1bb7c3-1f62-244c-0f3a-761cc1688042\",\"bb2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_id\":\"secondTraintupleID\"}]}"]}' -C myc
```
##### Command output:
```json
{
 "aggregatetuple_keys": null,
 "clean_models": false,
 "composite_traintuple_keys": null,
 "compute_plan_id": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "done_count": 0,
 "id_to_key": {
  "firstTraintupleID": "ab4840de60f351682c8060c5d151ce295f7c5432140da81d208d42b68657b7b2",
  "secondTraintupleID": "7e010308e046eaa6b73f41886fa57e6453032d6a3be7952d38fb66baaa61baca"
 },
 "metadata": {},
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "33f0b0c6656aabe113b47807492a1db295b76bc9d1d7d143f4bf3653014a2064"
 ],
 "traintuple_keys": [
  "ab4840de60f351682c8060c5d151ce295f7c5432140da81d208d42b68657b7b2",
  "7e010308e046eaa6b73f41886fa57e6453032d6a3be7952d38fb66baaa61baca"
 ],
 "tuple_count": 3
}
```
#### ------------ Update a ComputePlan ------------
Smart contract: `updateComputePlan`

##### JSON Inputs:
```go
{
 "compute_plan_id": string (required,required,len=64,hexadecimal),
 "traintuples": (omitempty) [{
   "data_manager_key": string (required,len=36),
   "data_sample_keys": [string] (required,dive,len=36),
   "algo_key": string (required,len=36),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 }],
 "aggregatetuples": (omitempty) [{
   "algo_key": string (required,len=36),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "worker": string (required),
 }],
 "composite_traintuples": (omitempty) [{
   "data_manager_key": string (required,len=36),
   "data_sample_keys": [string] (required,dive,len=36),
   "algo_key": string (required,len=36),
   "id": string (required,lte=64),
   "in_head_model_id": string (required_with=InTrunkModelID,omitempty,len=64,hexadecimal),
   "in_trunk_model_id": string (required_with=InHeadModelID,omitempty,len=64,hexadecimal),
   "out_trunk_model_permissions": (required){
     "process": (required){
       "public": bool (required),
       "authorized_ids": [string] (required),
     },
   },
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 }],
 "testtuples": (omitempty) [{
   "data_manager_key": string (omitempty,len=36),
   "data_sample_keys": [string] (omitempty,dive,len=36),
   "objective_key": string (required,len=36),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "traintuple_id": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateComputePlan","{\"compute_plan_id\":\"7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17\",\"traintuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c3-1f62-244c-0f3a-761cc1688042\"],\"algo_key\":\"fd1bb7c3-1f62-244c-0f3a-761cc1688042\",\"id\":\"thirdTraintupleID\",\"in_models_ids\":[\"firstTraintupleID\",\"secondTraintupleID\"],\"tag\":\"\",\"metadata\":null}],\"aggregatetuples\":null,\"composite_traintuples\":null,\"testtuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"bb1bb7c3-1f62-244c-0f3a-761cc1688042\",\"bb2bb7c3-1f62-244c-0f3a-761cc1688042\"],\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_id\":\"thirdTraintupleID\"}]}"]}' -C myc
```
##### Command output:
```json
{
 "aggregatetuple_keys": null,
 "clean_models": false,
 "composite_traintuple_keys": null,
 "compute_plan_id": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "done_count": 0,
 "id_to_key": {
  "thirdTraintupleID": "ceeb535f87810c425f7b23fde088f568ee2e174641f6bb3644d0eca74aa68937"
 },
 "metadata": {},
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "33f0b0c6656aabe113b47807492a1db295b76bc9d1d7d143f4bf3653014a2064",
  "7af1d164101e4991579a7dc13295f555123533659a4aa342acae099f58c952ba"
 ],
 "traintuple_keys": [
  "ab4840de60f351682c8060c5d151ce295f7c5432140da81d208d42b68657b7b2",
  "7e010308e046eaa6b73f41886fa57e6453032d6a3be7952d38fb66baaa61baca",
  "ceeb535f87810c425f7b23fde088f568ee2e174641f6bb3644d0eca74aa68937"
 ],
 "tuple_count": 5
}
```
#### ------------ Query an ObjectiveLeaderboard ------------
Smart contract: `queryObjectiveLeaderboard`

##### JSON Inputs:
```go
{
 "objective_key": string (omitempty,len=36),
 "ascendingOrder": bool (required),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryObjectiveLeaderboard","{\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"ascendingOrder\":true}"]}' -C myc
```
##### Command output:
```json
{
 "objective": {
  "description": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/description"
  },
  "key": "5c1d9cd1-c2c1-082d-de09-21b56d11030c",
  "metadata": {},
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "name": "accuracy",
   "storage_address": "https://toto/objective/222/metrics"
  },
  "name": "MSI classification",
  "owner": "SampleOrg",
  "permissions": {
   "process": {
    "authorized_ids": [],
    "public": true
   }
  },
  "test_dataset": {
   "data_manager_key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "data_sample_keys": [
    "bb1bb7c3-1f62-244c-0f3a-761cc1688042",
    "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
   ],
   "metadata": {},
   "worker": ""
  }
 },
 "testtuples": [
  {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storage_address": "https://toto/algo/222/algo"
   },
   "creator": "SampleOrg",
   "key": "f908806c4341efa2d5fef70e0b3171017f5315dc31ce1383c7bf5585653ef883",
   "perf": 0.9,
   "tag": "",
   "traintuple_key": "b0289ab83a71f01e2b720259a645224453e841ff0c3335b874b61c33344f8a75"
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
 "aggregatetuple_keys": null,
 "clean_models": false,
 "composite_traintuple_keys": null,
 "compute_plan_id": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "done_count": 0,
 "id_to_key": {},
 "metadata": {},
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "33f0b0c6656aabe113b47807492a1db295b76bc9d1d7d143f4bf3653014a2064",
  "7af1d164101e4991579a7dc13295f555123533659a4aa342acae099f58c952ba"
 ],
 "traintuple_keys": [
  "ab4840de60f351682c8060c5d151ce295f7c5432140da81d208d42b68657b7b2",
  "7e010308e046eaa6b73f41886fa57e6453032d6a3be7952d38fb66baaa61baca",
  "ceeb535f87810c425f7b23fde088f568ee2e174641f6bb3644d0eca74aa68937"
 ],
 "tuple_count": 5
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
  "aggregatetuple_keys": null,
  "clean_models": false,
  "composite_traintuple_keys": null,
  "compute_plan_id": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
  "done_count": 0,
  "id_to_key": {},
  "metadata": {},
  "status": "todo",
  "tag": "a tag is simply a string",
  "testtuple_keys": [
   "33f0b0c6656aabe113b47807492a1db295b76bc9d1d7d143f4bf3653014a2064",
   "7af1d164101e4991579a7dc13295f555123533659a4aa342acae099f58c952ba"
  ],
  "traintuple_keys": [
   "ab4840de60f351682c8060c5d151ce295f7c5432140da81d208d42b68657b7b2",
   "7e010308e046eaa6b73f41886fa57e6453032d6a3be7952d38fb66baaa61baca",
   "ceeb535f87810c425f7b23fde088f568ee2e174641f6bb3644d0eca74aa68937"
  ],
  "tuple_count": 5
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
 "aggregatetuple_keys": null,
 "clean_models": false,
 "composite_traintuple_keys": null,
 "compute_plan_id": "7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17",
 "done_count": 0,
 "id_to_key": {},
 "metadata": {},
 "status": "canceled",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "33f0b0c6656aabe113b47807492a1db295b76bc9d1d7d143f4bf3653014a2064",
  "7af1d164101e4991579a7dc13295f555123533659a4aa342acae099f58c952ba"
 ],
 "traintuple_keys": [
  "ab4840de60f351682c8060c5d151ce295f7c5432140da81d208d42b68657b7b2",
  "7e010308e046eaa6b73f41886fa57e6453032d6a3be7952d38fb66baaa61baca",
  "ceeb535f87810c425f7b23fde088f568ee2e174641f6bb3644d0eca74aa68937"
 ],
 "tuple_count": 5
}
```
