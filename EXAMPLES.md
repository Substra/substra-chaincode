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
 "name": string (required,gte=1,lte=100),
 "opener_hash": string (required,len=64,hexadecimal),
 "opener_storage_address": string (required,url),
 "type": string (required,gte=1,lte=30),
 "description_hash": string (required,len=64,hexadecimal),
 "description_storage_address": string (required,url),
 "objective_key": string (omitempty),
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
peer chaincode invoke -n mycc -c '{"Args":["registerDataManager","{\"name\":\"liver slide\",\"opener_hash\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"opener_storage_address\":\"https://toto/dataManager/42234/opener\",\"type\":\"images\",\"description_hash\":\"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee\",\"description_storage_address\":\"https://toto/dataManager/42234/description\",\"objective_key\":\"\",\"permissions\":{\"process\":{\"public\":true,\"authorized_ids\":[]}},\"metadata\":null}"]}' -C myc
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
  "storage_address": "https://toto/dataManager/42234/description"
 },
 "key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
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
 "data_manager_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"hashes\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"data_manager_keys\":[\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"testOnly\":\"true\"}"]}' -C myc
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
 "description_hash": string (required,len=64,hexadecimal),
 "description_storage_address": string (required,url),
 "metrics_name": string (required,gte=1,lte=100),
 "metrics_hash": string (required,len=64,hexadecimal),
 "metrics_storage_address": string (required,url),
 "test_dataset": (omitempty){
   "data_manager_key": string (omitempty,len=64,hexadecimal),
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
peer chaincode invoke -n mycc -c '{"Args":["registerObjective","{\"name\":\"MSI classification\",\"description_hash\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"description_storage_address\":\"https://toto/objective/222/description\",\"metrics_name\":\"accuracy\",\"metrics_hash\":\"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"metrics_storage_address\":\"https://toto/objective/222/metrics\",\"test_dataset\":{\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"]},\"permissions\":{\"process\":{\"public\":true,\"authorized_ids\":[]}},\"metadata\":null}"]}' -C myc
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
 "data_manager_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "testOnly": string (required,oneof=true false),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample","{\"hashes\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"data_manager_keys\":[\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"testOnly\":\"false\"}"]}' -C myc
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
  "key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "metadata": {},
  "name": "liver slide",
  "objective_key": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
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
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
   "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "key": "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "owner": "SampleOrg"
 },
 {
  "data_manager_keys": [
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
   "storage_address": "https://toto/objective/222/description"
  },
  "key": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
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
   "data_manager_key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
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
 "data_manager_key": string (required,len=64,hexadecimal),
 "data_sample_keys": [string] (required,unique,gt=0,dive,len=64,hexadecimal),
 "compute_plan_id": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"in_models\":[],\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"compute_plan_id\":\"\",\"rank\":\"\",\"tag\":\"\",\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c"
}
```
#### ------------ Add Traintuple with inModel from previous traintuple ------------
Smart contract: `createTraintuple`

##### JSON Inputs:
```go
{
 "algo_key": string (required,len=64,hexadecimal),
 "in_models": [string] (omitempty,dive,len=64,hexadecimal),
 "data_manager_key": string (required,len=64,hexadecimal),
 "data_sample_keys": [string] (required,unique,gt=0,dive,len=64,hexadecimal),
 "compute_plan_id": string (omitempty),
 "rank": string (omitempty),
 "tag": string (omitempty,lte=64),
 "metadata": map (lte=100,dive,keys,lte=50,endkeys,lte=100),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","{\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"in_models\":[\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"],\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"compute_plan_id\":\"\",\"rank\":\"\",\"tag\":\"\",\"metadata\":null}"]}' -C myc
```
##### Command output:
```json
{
 "key": "ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d"
}
```
##### Command output:
```json
{
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
   "storage_address": "https://toto/algo/222/algo"
  },
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "metadata": {},
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "in_models": null,
  "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
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
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\",\"log\":\"no error, ah ah ah\",\"out_model\":{\"hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"storage_address\":\"https://substrabac/model/toto\"}}"]}' -C myc
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
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
peer chaincode invoke -n mycc -c '{"Args":["queryTraintuple","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
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
  "keys": [
   "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "metadata": {},
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "worker": "SampleOrg"
 },
 "in_models": null,
 "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
 "data_manager_key": string (omitempty,len=64,hexadecimal),
 "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "objective_key": string (required,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objective_key\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8"
}
```
#### ------------ Add Certified Testtuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "data_manager_key": string (omitempty,len=64,hexadecimal),
 "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "objective_key": string (required,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"\",\"data_sample_keys\":null,\"objective_key\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
```
##### Command output:
```json
{
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d"
}
```
#### ------------ Add Testtuple with not done traintuple ------------
Smart contract: `createTesttuple`

##### JSON Inputs:
```go
{
 "data_manager_key": string (omitempty,len=64,hexadecimal),
 "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
 "objective_key": string (required,len=64,hexadecimal),
 "tag": string (omitempty,lte=64),
 "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
 "traintuple_key": string (required,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","{\"data_manager_key\":\"\",\"data_sample_keys\":null,\"objective_key\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"metadata\":null,\"traintuple_key\":\"ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d\"}"]}' -C myc
```
##### Command output:
```json
{
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
   "storage_address": "https://toto/algo/222/algo"
  },
  "certified": false,
  "compute_plan_id": "",
  "creator": "SampleOrg",
  "dataset": {
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","{\"key\":\"a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d\"}"]}' -C myc
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
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0,
  "worker": "SampleOrg"
 },
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
 "log": "",
 "metadata": {},
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "doing",
 "tag": "",
 "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","{\"key\":\"a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d\",\"log\":\"no error, ah ah ah\",\"perf\":0.9}"]}' -C myc
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
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
 "log": "no error, ah ah ah",
 "metadata": {},
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "done",
 "tag": "",
 "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
peer chaincode query -n mycc -c '{"Args":["queryTesttuple","{\"key\":\"a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d\"}"]}' -C myc
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
  "keys": [
   "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
  ],
  "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
  "perf": 0.9,
  "worker": "SampleOrg"
 },
 "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
 "log": "no error, ah ah ah",
 "metadata": {},
 "objective": {
  "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
  "metrics": {
   "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/metrics"
  }
 },
 "rank": 0,
 "status": "done",
 "tag": "",
 "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "todo",
  "tag": "",
  "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
  "log": "no error, ah ah ah",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "done",
  "tag": "",
  "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0,
   "worker": "SampleOrg"
  },
  "key": "f5515eef9906b6af129355146a648c94f390e32d0fb45a8f54fddfe3329df716",
  "log": "",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "waiting",
  "tag": "",
  "traintuple_key": "ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d",
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
peer chaincode query -n mycc -c '{"Args":["queryModelDetails","{\"key\":\"ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c\"}"]}' -C myc
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
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "perf": 0,
    "worker": "SampleOrg"
   },
   "key": "6ad32c063f8f1ae04626987e0b15351c3a2007a417ba2bdc557b7ff4c7a9ebf8",
   "log": "",
   "metadata": {},
   "objective": {
    "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "metrics": {
     "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
     "storage_address": "https://toto/objective/222/metrics"
    }
   },
   "rank": 0,
   "status": "todo",
   "tag": "",
   "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
   "keys": [
    "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "perf": 0.9,
   "worker": "SampleOrg"
  },
  "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
  "log": "no error, ah ah ah",
  "metadata": {},
  "objective": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "metrics": {
    "hash": "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
    "storage_address": "https://toto/objective/222/metrics"
   }
  },
  "rank": 0,
  "status": "done",
  "tag": "",
  "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
   "keys": [
    "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
   ],
   "metadata": {},
   "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
   "worker": "SampleOrg"
  },
  "in_models": null,
  "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
    "keys": [
     "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
     "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
    ],
    "metadata": {},
    "opener_hash": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
    "worker": "SampleOrg"
   },
   "in_models": null,
   "key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c",
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
     "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c"
    }
   ],
   "key": "ed8102d4f4e19e961585a0b544c76c87c9ffeaf1bcbec57247023e240e3bde2d",
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
  "storage_address": "https://toto/dataManager/42234/description"
 },
 "key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
 "metadata": {},
 "name": "liver slide",
 "objective_key": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
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
 "data_manager_keys": [string] (required,dive,len=64,hexadecimal),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateDataSample","{\"hashes\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"data_manager_keys\":[\"38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee\"]}"]}' -C myc
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
  "storage_address": "https://toto/dataManager/42234/description"
 },
 "key": "38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee",
 "metadata": {},
 "name": "liver slide",
 "objective_key": "",
 "opener": {
  "hash": "38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee",
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
   "data_manager_key": string (required,len=64,hexadecimal),
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
   "data_manager_key": string (required,len=64,hexadecimal),
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
   "data_manager_key": string (omitempty,len=64,hexadecimal),
   "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
   "objective_key": string (required,len=64,hexadecimal),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "traintuple_id": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["createComputePlan","{\"clean_models\":false,\"tag\":\"a tag is simply a string\",\"metadata\":null,\"traintuples\":[{\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"firstTraintupleID\",\"in_models_ids\":null,\"tag\":\"\",\"metadata\":null},{\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"secondTraintupleID\",\"in_models_ids\":[\"firstTraintupleID\"],\"tag\":\"\",\"metadata\":null}],\"aggregatetuples\":null,\"composite_traintuples\":null,\"testtuples\":[{\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objective_key\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"metadata\":null,\"traintuple_id\":\"secondTraintupleID\"}]}"]}' -C myc
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
  "firstTraintupleID": "01feb56691d26983a641d29f4c2a5b7098f99eb471b7e5f03aaa78c8ae142ca9",
  "secondTraintupleID": "78914b1f480f5e81a26e4d04d88bdb27937e858c49c6bb9d1ae83ff6627ca0ab"
 },
 "metadata": {},
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "1ca3227c1a1232a55e31d11d93b1fe224f454c5a4508093a15a6cae2a220f957"
 ],
 "traintuple_keys": [
  "01feb56691d26983a641d29f4c2a5b7098f99eb471b7e5f03aaa78c8ae142ca9",
  "78914b1f480f5e81a26e4d04d88bdb27937e858c49c6bb9d1ae83ff6627ca0ab"
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
   "data_manager_key": string (required,len=64,hexadecimal),
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
   "data_manager_key": string (required,len=64,hexadecimal),
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
   "data_manager_key": string (omitempty,len=64,hexadecimal),
   "data_sample_keys": [string] (omitempty,dive,len=64,hexadecimal),
   "objective_key": string (required,len=64,hexadecimal),
   "tag": string (omitempty,lte=64),
   "metadata": map (omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100),
   "traintuple_id": string (required,lte=64),
 }],
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["updateComputePlan","{\"compute_plan_id\":\"7dd808239c1e062399449bd11b634d9bd1fd0a2b795ad345b62f95b4933bfa17\",\"traintuples\":[{\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"algo_key\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"id\":\"thirdTraintupleID\",\"in_models_ids\":[\"firstTraintupleID\",\"secondTraintupleID\"],\"tag\":\"\",\"metadata\":null}],\"aggregatetuples\":null,\"composite_traintuples\":null,\"testtuples\":[{\"data_manager_key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"data_sample_keys\":[\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"],\"objective_key\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"tag\":\"\",\"metadata\":null,\"traintuple_id\":\"thirdTraintupleID\"}]}"]}' -C myc
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
  "thirdTraintupleID": "17c7623e87be77d8f93f21401e2eae98384de4a1d7ee841c0b9e0a07897cfbbf"
 },
 "metadata": {},
 "status": "todo",
 "tag": "a tag is simply a string",
 "testtuple_keys": [
  "1ca3227c1a1232a55e31d11d93b1fe224f454c5a4508093a15a6cae2a220f957",
  "4395e03f727aa5ce5c3d8eb034e23871e57d9a6546a57b9a6f786bf1019e0b52"
 ],
 "traintuple_keys": [
  "01feb56691d26983a641d29f4c2a5b7098f99eb471b7e5f03aaa78c8ae142ca9",
  "78914b1f480f5e81a26e4d04d88bdb27937e858c49c6bb9d1ae83ff6627ca0ab",
  "17c7623e87be77d8f93f21401e2eae98384de4a1d7ee841c0b9e0a07897cfbbf"
 ],
 "tuple_count": 5
}
```
#### ------------ Query an ObjectiveLeaderboard ------------
Smart contract: `queryObjectiveLeaderboard`

##### JSON Inputs:
```go
{
 "objective_key": string (omitempty,len=64,hexadecimal),
 "ascendingOrder": bool (required),
}
```
##### Command peer example:
```bash
peer chaincode invoke -n mycc -c '{"Args":["queryObjectiveLeaderboard","{\"objective_key\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"ascendingOrder\":true}"]}' -C myc
```
##### Command output:
```json
{
 "objective": {
  "description": {
   "hash": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
   "storage_address": "https://toto/objective/222/description"
  },
  "key": "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
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
   "data_manager_key": "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
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
   "key": "a0f368a23449ae1751ccb2335f79d8ff084bc7bb13e1e2b5252d930857bc4d2d",
   "perf": 0.9,
   "tag": "",
   "traintuple_key": "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c"
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
  "1ca3227c1a1232a55e31d11d93b1fe224f454c5a4508093a15a6cae2a220f957",
  "4395e03f727aa5ce5c3d8eb034e23871e57d9a6546a57b9a6f786bf1019e0b52"
 ],
 "traintuple_keys": [
  "01feb56691d26983a641d29f4c2a5b7098f99eb471b7e5f03aaa78c8ae142ca9",
  "78914b1f480f5e81a26e4d04d88bdb27937e858c49c6bb9d1ae83ff6627ca0ab",
  "17c7623e87be77d8f93f21401e2eae98384de4a1d7ee841c0b9e0a07897cfbbf"
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
   "1ca3227c1a1232a55e31d11d93b1fe224f454c5a4508093a15a6cae2a220f957",
   "4395e03f727aa5ce5c3d8eb034e23871e57d9a6546a57b9a6f786bf1019e0b52"
  ],
  "traintuple_keys": [
   "01feb56691d26983a641d29f4c2a5b7098f99eb471b7e5f03aaa78c8ae142ca9",
   "78914b1f480f5e81a26e4d04d88bdb27937e858c49c6bb9d1ae83ff6627ca0ab",
   "17c7623e87be77d8f93f21401e2eae98384de4a1d7ee841c0b9e0a07897cfbbf"
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
  "1ca3227c1a1232a55e31d11d93b1fe224f454c5a4508093a15a6cae2a220f957",
  "4395e03f727aa5ce5c3d8eb034e23871e57d9a6546a57b9a6f786bf1019e0b52"
 ],
 "traintuple_keys": [
  "01feb56691d26983a641d29f4c2a5b7098f99eb471b7e5f03aaa78c8ae142ca9",
  "78914b1f480f5e81a26e4d04d88bdb27937e858c49c6bb9d1ae83ff6627ca0ab",
  "17c7623e87be77d8f93f21401e2eae98384de4a1d7ee841c0b9e0a07897cfbbf"
 ],
 "tuple_count": 5
}
```
