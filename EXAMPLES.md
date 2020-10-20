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
 "key": string (required,len=64,hexadecimal),
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
 "hashes": [string] (required,dive,len=64,hexadecimal),
 "data_manager_keys": [string] (omitempty,dive,len=36),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"hashes\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"data_manager_keys\":[\"da1bb7c3-1f62-244c-0f3a-761cc1688042\"],\"testOnly\":\"true\"}"]}' -C myc
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
 "key": string (required,len=36),
 "name": string (required,gte=1,lte=100),
 "description_hash": string (required,len=64,hexadecimal),
 "description_storage_address": string (required,url),
 "metrics_name": string (required,gte=1,lte=100),
 "metrics_hash": string (required,len=64,hexadecimal),
 "metrics_storage_address": string (required,url),
 "test_dataset": (omitempty){
   "data_manager_key": string (omitempty,len=36),
   "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
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
peer chaincode invoke -n mycc -c '{"Args":["registerObjective","{\"key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"name\":\"MSI classification\",\"description_hash\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"description_storage_address\":\"https://toto/objective/222/description\",\"metrics_name\":\"accuracy\",\"metrics_hash\":\"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"metrics_storage_address\":\"https://toto/objective/222/metrics\",\"test_dataset\":{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"]},\"permissions\":{\"process\":{\"public\":true,\"authorized_ids\":[]}},\"metadata\":null}"]}' -C myc
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
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo","{\"name\":\"hog + svm\",\"hash\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"storage_address\":\"https://toto/algo/222/algo\",\"description_hash\":\"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca\",\"description_storage_address\":\"https://toto/algo/222/description\",\"permissions\":{\"process\":{\"public\":true,\"authorized_ids\":[]}},\"metadata\":null}"]}' -C myc
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
 "data_manager_keys": [string] (omitempty,dive,len=36),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"hashes\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"data_manager_keys\":[\"da1bb7c3-1f62-244c-0f3a-761cc1688042\"],\"testOnly\":\"false\"}"]}' -C myc
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
  "key": "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "key": "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c3-1f62-244c-0f3a-761cc1688042"
  ],
  "key": "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c3-1f62-244c-0f3a-761cc1688042"
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
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
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
 "algo_key": string (required,len=64,hexadecimal),
 "in_models": [string] (omitempty,dive,len=64,hexadecimal),
 "data_manager_key": string (required,len=36),
 "data_sample_keys": [string] (required,unique,gt=0,dive,len=64,hexadecimal),
 "compute_plan_id": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"in_models\":[],\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"compute_plan_id\":\"\",\"rank\":\"\",\"tag\":\"\",\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765"
}
```
#### ------------ Add Traintuple with inModel from previous traintuple ------------
Smart contract: `createTraintuple`

##### JSON Inputs:
```go
{
 "algo_key": string (required,len=64,hexadecimal),
 "in_models": [string] (omitempty,dive,len=64,hexadecimal),
 "data_manager_key": string (required,len=36),
 "data_sample_keys": [string] (required,unique,gt=0,dive,len=64,hexadecimal),
 "compute_plan_id": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"in_models\":[\"5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765\"],\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"compute_plan_id\":\"\",\"rank\":\"\",\"tag\":\"\",\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "6061e82631b923157fa53e6d22fa02ebfa69d731734aa9099396e46a297d5ec5"
}
```
##### Command output:
```json
{
 "key": "6061e82631b923157fa53e6d22fa02ebfa69d731734aa9099396e46a297d5ec5"
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
   "storage_address": "https://toto/algo/222/algo"
  },
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "metadata": {},
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "in_models": null,
  "key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain","{\"key\":\"5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","{\"key\":\"5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765\",\"log\":\"no error, ah ah ah\",\"out_model\":{\"hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"storage_address\":\"https://substrabac/model/toto\"}}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
peer chaincode invoke -n mycc -c '{"Args":["queryTraintuple","{\"key\":\"5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
 "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "objective_key": string (required,len=36),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "7ca96c767379bda1046eb870f216c7881afa461d78cfc24f4ac9441f27c723d8"
}
```
#### ------------ Add Certified Testtuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "data_manager_key": string (omitempty,len=36),
 "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "objective_key": string (required,len=36),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"\",\"data_sample_keys\":null,\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f"
}
```
#### ------------ Add Testtuple with not done traintuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "data_manager_key": string (omitempty,len=36),
 "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "objective_key": string (required,len=36),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"\",\"data_sample_keys\":null,\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"6061e82631b923157fa53e6d22fa02ebfa69d731734aa9099396e46a297d5ec5\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "32507b460b0d0de26118fc16d05c58d93342b73c0250ccf2d702c059c8eacd92"
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
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": false,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "7ca96c767379bda1046eb870f216c7881afa461d78cfc24f4ac9441f27c723d8",
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
  "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
  "traintuple_type": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f",
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
  "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","{\"key\":\"e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "certified": true,
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0,
  "worker": "SampleOrg"
 },
 "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f",
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
 "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","{\"key\":\"e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f\",\"log\":\"no error, ah ah ah\",\"perf\":0.9}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "certified": true,
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f",
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
 "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
peer chaincode query -n mycc -c '{"Args":["queryTesttuple","{\"key\":\"e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f\"}"]}' -C myc
```
##### Command output:
```json
{
 "algo": {
  "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "name": "hog + svm",
  "storage_address": "https://toto/algo/222/algo"
 },
 "certified": true,
 "compute_plan_id": "",
 "creator": "SampleOrg",
 "dataset": {
  "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f",
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
 "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": false,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "7ca96c767379bda1046eb870f216c7881afa461d78cfc24f4ac9441f27c723d8",
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
  "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
  "traintuple_type": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f",
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
  "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
  "traintuple_type": "traintuple"
 },
 {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "32507b460b0d0de26118fc16d05c58d93342b73c0250ccf2d702c059c8eacd92",
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
  "traintuple_key": "6061e82631b923157fa53e6d22fa02ebfa69d731734aa9099396e46a297d5ec5",
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
peer chaincode query -n mycc -c '{"Args":["queryModelDetails","{\"key\":\"5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765\"}"]}' -C myc
```
##### Command output:
```json
{
 "non_certified_testtuples": [
  {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storage_address": "https://toto/algo/222/algo"
   },
   "certified": false,
   "compute_plan_id": "",
   "creator": "SampleOrg",
   "dataset": {
    "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0,
    "worker": "SampleOrg"
   },
   "key": "7ca96c767379bda1046eb870f216c7881afa461d78cfc24f4ac9441f27c723d8",
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
   "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
   "traintuple_type": "traintuple"
  }
 ],
 "testtuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": true,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f",
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
  "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
  "traintuple_type": "traintuple"
 },
 "traintuple": {
  "algo": {
   "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "name": "hog + svm",
   "storage_address": "https://toto/algo/222/algo"
  },
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "metadata": {},
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "in_models": null,
  "key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
    "name": "hog + svm",
    "storage_address": "https://toto/algo/222/algo"
   },
   "compute_plan_id": "",
   "creator": "SampleOrg",
   "dataset": {
    "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "metadata": {},
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "worker": "SampleOrg"
   },
   "in_models": null,
   "key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765",
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
 },
 {
  "traintuple": {
   "algo": {
    "hash": "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "name": "hog + svm",
    "storage_address": "https://toto/algo/222/algo"
   },
   "compute_plan_id": "",
   "creator": "SampleOrg",
   "dataset": {
    "key": "da1bb7c3-1f62-244c-0f3a-761cc1688042",
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "metadata": {},
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "worker": "SampleOrg"
   },
   "in_models": [
    {
     "hash": "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed",
     "storage_address": "https://substrabac/model/toto",
     "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765"
    }
   ],
   "key": "6061e82631b923157fa53e6d22fa02ebfa69d731734aa9099396e46a297d5ec5",
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
  "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
 ],
 "train_data_sample_keys": [
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
 "data_manager_keys": [string] (required,dive,len=36),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateDataSample","{\"hashes\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"data_manager_keys\":[\"38a320b2-a67c-8003-cc74-8d6666534f2b\"]}"]}' -C myc
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
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuples": (omitempty) [{
   "data_manager_key": string (required,len=36),
   "data_sample_keys": [string] (required,dive,len=64,hexadecimal),
   "algo_key": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 }],
 "aggregatetuples": (omitempty) [{
   "algo_key": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "worker": string (required),
 }],
 "composite_traintuples": (omitempty) [{
   "data_manager_key": string (required,len=36),
   "data_sample_keys": [string] (required,dive,len=64,hexadecimal),
   "algo_key": string (required,len=64,hexadecimal),
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
   "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
   "objective_key": string (required,len=36),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "traintuple_id": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createComputePlan","{\"clean_models\":false,\"tag\":\"a tag is simply a string\",\"metadata\":null,\"traintuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"firstTraintupleID\",\"in_models_ids\":null,\"tag\":\"\",\"metadata\":null},{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"secondTraintupleID\",\"in_models_ids\":[\"firstTraintupleID\"],\"tag\":\"\",\"metadata\":null}],\"aggregatetuples\":null,\"composite_traintuples\":null,\"testtuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_id\":\"secondTraintupleID\"}]}"]}' -C myc
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
  "firstTraintupleID": "47ec9d5d518277c294167d716a91de1a25742f0bbcf41d2f52a4c77c1026fa8a",
  "secondTraintupleID": "07e4fc257283ec86a88c845368f7f2010509c65fb2495b91d3c5530bb101d45a"
 },
 "metadata": {},
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "e13484a0c7eaf152ec1c6abd00a5c1609a5fd0e605eb2d2cc66108e5a4790aed"
 ],
 "traintuple_keys": [
  "47ec9d5d518277c294167d716a91de1a25742f0bbcf41d2f52a4c77c1026fa8a",
  "07e4fc257283ec86a88c845368f7f2010509c65fb2495b91d3c5530bb101d45a"
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
   "data_sample_keys": [string] (required,dive,len=64,hexadecimal),
   "algo_key": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 }],
 "aggregatetuples": (omitempty) [{
   "algo_key": string (required,len=64,hexadecimal),
   "id": string (required,lte=64),
   "in_models_ids": [string] (omitempty,dive,lte=64),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "worker": string (required),
 }],
 "composite_traintuples": (omitempty) [{
   "data_manager_key": string (required,len=36),
   "data_sample_keys": [string] (required,dive,len=64,hexadecimal),
   "algo_key": string (required,len=64,hexadecimal),
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
   "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
   "objective_key": string (required,len=36),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "traintuple_id": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateComputePlan","{\"compute_plan_id\":\"7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17\",\"traintuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"thirdTraintupleID\",\"in_models_ids\":[\"firstTraintupleID\",\"secondTraintupleID\"],\"tag\":\"\",\"metadata\":null}],\"aggregatetuples\":null,\"composite_traintuples\":null,\"testtuples\":[{\"data_manager_key\":\"da1bb7c3-1f62-244c-0f3a-761cc1688042\",\"data_sample_keys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objective_key\":\"5c1d9cd1-c2c1-082d-de09-21b56d11030c\",\"tag\":\"\",\"metadata\":null,\"traintuple_id\":\"thirdTraintupleID\"}]}"]}' -C myc
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
  "thirdTraintupleID": "f0f06e6d772d7ab4dd67cab62c2bcacb04b3957d4714c124da02d7b9ad386639"
 },
 "metadata": {},
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "e13484a0c7eaf152ec1c6abd00a5c1609a5fd0e605eb2d2cc66108e5a4790aed",
  "0750df946a10c96f0285ef92da097b364d798c2505a4e398e4bc49280cd3200b"
 ],
 "traintuple_keys": [
  "47ec9d5d518277c294167d716a91de1a25742f0bbcf41d2f52a4c77c1026fa8a",
  "07e4fc257283ec86a88c845368f7f2010509c65fb2495b91d3c5530bb101d45a",
  "f0f06e6d772d7ab4dd67cab62c2bcacb04b3957d4714c124da02d7b9ad386639"
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
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
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
   "key": "e59d1804713fcf0badcb82a82ca1df0e3094af02706ec44e544bccdcedbe1d6f",
   "perf": 0.9,
   "tag": "",
   "traintuple_key": "5ec4813058be788b991baa236da13f875b82de6add27345ae6049ee820976765"
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
  "e13484a0c7eaf152ec1c6abd00a5c1609a5fd0e605eb2d2cc66108e5a4790aed",
  "0750df946a10c96f0285ef92da097b364d798c2505a4e398e4bc49280cd3200b"
 ],
 "traintuple_keys": [
  "47ec9d5d518277c294167d716a91de1a25742f0bbcf41d2f52a4c77c1026fa8a",
  "07e4fc257283ec86a88c845368f7f2010509c65fb2495b91d3c5530bb101d45a",
  "f0f06e6d772d7ab4dd67cab62c2bcacb04b3957d4714c124da02d7b9ad386639"
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
   "e13484a0c7eaf152ec1c6abd00a5c1609a5fd0e605eb2d2cc66108e5a4790aed",
   "0750df946a10c96f0285ef92da097b364d798c2505a4e398e4bc49280cd3200b"
  ],
  "traintuple_keys": [
   "47ec9d5d518277c294167d716a91de1a25742f0bbcf41d2f52a4c77c1026fa8a",
   "07e4fc257283ec86a88c845368f7f2010509c65fb2495b91d3c5530bb101d45a",
   "f0f06e6d772d7ab4dd67cab62c2bcacb04b3957d4714c124da02d7b9ad386639"
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
  "e13484a0c7eaf152ec1c6abd00a5c1609a5fd0e605eb2d2cc66108e5a4790aed",
  "0750df946a10c96f0285ef92da097b364d798c2505a4e398e4bc49280cd3200b"
 ],
 "traintuple_keys": [
  "47ec9d5d518277c294167d716a91de1a25742f0bbcf41d2f52a4c77c1026fa8a",
  "07e4fc257283ec86a88c845368f7f2010509c65fb2495b91d3c5530bb101d45a",
  "f0f06e6d772d7ab4dd67cab62c2bcacb04b3957d4714c124da02d7b9ad386639"
 ],
 "tuple_count": 5
}
```
