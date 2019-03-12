# substra-chaincode
Chaincode for the Substra platform

> :warning: This project is under active development. Please, wait some times before using it...

## License

This project is developed under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](./LICENSE) file.
## Devmode

See [chaincode-docker-devmode](./chaincode-docker-devmode/README.rst)

## Documentation


Note for internal use only: See the [technical specifications](https://github.com/SubstraFoundation/substra-spec/blob/master/technical/technical_spec_substra.md#p0_smartcontract).

### Implemented smart contracts


- `registerProblem`
- `registerDataset`
- `updateDataset`
- `registerData`
- `updateData`
- `registerAlgo`
- `createTraintuple`
- `logStartTrain`
- `logSuccessTrain`
- `logFailTrain`
- `createTesttuple`
- `logStartTest`
- `logSuccessTest`
- `logFailTest`
- `query`
- `queryProblems`
- `queryAlgo`
- `queryTraintuples`
- `queryDatasets`
- `queryModelDetails`
- `queryModels`
- `queryDatasetData`
- `queryFilter`

### Examples 
#### ------------ Add Dataset ------------
Smart contract: `registerDataset`  
 Inputs: `Name`, `OpenerHash`, `OpenerStorageAddress`, `Type`, `DescriptionHash`, `DescriptionStorageAddress`, `ChallengeKey`, `Permissions`
```
peer chaincode invoke -n mycc -c '{"Args":["registerDataset","liver slide","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","https://toto/dataset/42234/opener","images","8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","https://toto/dataset/42234/description","","all"]}' -C myc
```
>  da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc 

#### ------------ Query Dataset From key ------------
Smart contract: `query`  
 Inputs: `elementKey`
```
peer chaincode query -n mycc -c '{"Args":["query","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  {"name":"liver slide","openerStorageAddress":"https://toto/dataset/42234/opener","type":"images","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","challengeKey":"","permissions":"all"} 

#### ------------ Add test Data ------------
Smart contract: `registerData`  
 Inputs: `Hashes`, `DatasetKeys`, `TestOnly`
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","true"]}' -C myc
```
>  {"keys": ["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc", "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]} 

#### ------------ Add Challenge ------------
Smart contract: `registerChallenge`  
 Inputs: `Name`, `DescriptionHash`, `DescriptionStorageAddress`, `MetricsName`, `MetricsHash`, `MetricsStorageAddress`, `TestData`, `Permissions`
```
peer chaincode invoke -n mycc -c '{"Args":["registerChallenge","MSI classification","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","https://toto/challenge/222/description","accuracy","4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","https://toto/challenge/222/metrics","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","all"]}' -C myc
```
>  5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379 

#### ------------ Add Algo ------------
Smart contract: `registerAlgo`  
 Inputs: `Name`, `Hash`, `StorageAddress`, `DescriptionHash`, `DescriptionStorageAddress`, `ChallengeKey`, `Permissions`
```
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo","hog + svm","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","https://toto/algo/222/algo","e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca","https://toto/algo/222/description","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","all"]}' -C myc
```
>  fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc 

#### ------------ Add Train Data ------------
Smart contract: `registerData`  
 Inputs: `Hashes`, `DatasetKeys`, `TestOnly`
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","false"]}' -C myc
```
>  {"keys": ["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc", "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]} 

#### ------------ Query Datasets ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasets"]}' -C myc
```
>  [{"challengeKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","openerStorageAddress":"https://toto/dataset/42234/opener","owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","type":"images"}] 

#### ------------ Query Challenges ------------
```
peer chaincode query -n mycc -c '{"Args":["queryChallenges"]}' -C myc
```
>  [{"descriptionStorageAddress":"https://toto/challenge/222/description","key":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","name":"accuracy","storageAddress":"https://toto/challenge/222/metrics"},"name":"MSI classification","owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","testData":{"dataKeys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"datasetKey":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"}}] 

#### ------------ Add Traintuple ------------
Smart contract: `createTraintuple`  
 Inputs: `AlgoKey`, `InModels`, `DatasetKey`, `DataKeys`, `FLtask`, `Rank`
```
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","",""]}' -C myc
```
>  337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f 

#### ------------ Add Traintuple with inModel from previous traintuple ------------
Smart contract: `createTraintuple`  
 Inputs: `AlgoKey`, `InModels`, `DatasetKey`, `DataKeys`, `FLtask`, `Rank`
```
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","",""]}' -C myc
```
>  f952b8514669261ca1c53ac853c3abab4a870fab395215ed8a9e7e308eb00c6b 

#### ------------ Query Traintuples of worker with todo status ------------
```
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","traintuple~worker~status",", todo"]}' -C myc
```
>  null 

#### ------------ Log Start Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain","337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"]}' -C myc
```
>  {"algoKey":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","inModels":null,"outModel":null,"data":{"datasetKey":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]},"perf":0,"fltask":"","rank":0,"status":"doing","log":"","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Log Success Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed, https://substrabac/model/toto","0.9","no error, ah ah ah"]}' -C myc
```
>  {"algoKey":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","inModels":null,"outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"data":{"datasetKey":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]},"perf":0.9,"fltask":"","rank":0,"status":"done","log":"no error, ah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Query Traintuple From key ------------
```
peer chaincode query -n mycc -c '{"Args":["query","337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"]}' -C myc
```
>  {"algoKey":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","inModels":null,"outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"data":{"datasetKey":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]},"perf":0.9,"fltask":"","rank":0,"status":"done","log":"no error, ah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Add Non-Certified Testtuple ------------
Smart contract: `createTesttuple`  
 Inputs: `TraintupleKey`, `DatasetKey`, `DataKeys`
```
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  5608124982f25ea50a5655da1ad91eace89142e79562cbf78b6f79c5d667c9c1 

#### ------------ Add Certified Testtuple ------------
Smart contract: `createTesttuple`  
 Inputs: `TraintupleKey`, `DatasetKey`, `DataKeys`
```
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","",""]}' -C myc
```
>  3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6 

#### ------------ Add Testtuple with not done traintuple ------------
Smart contract: `createTesttuple`  
 Inputs: `TraintupleKey`, `DatasetKey`, `DataKeys`
```
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple","f952b8514669261ca1c53ac853c3abab4a870fab395215ed8a9e7e308eb00c6b","",""]}' -C myc
```
>  671d7c25e98b195dbe99644bcff39d9f01f0fa3b2fc09c88f8453b9e965db0f2 

#### ------------ Query Testtuples of worker with todo status ------------
```
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","testtuple~worker~status","bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0, todo"]}' -C myc
```
>  [{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6","log":"","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"},"permissions":"all","status":"todo"},{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":false,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"5608124982f25ea50a5655da1ad91eace89142e79562cbf78b6f79c5d667c9c1","log":"","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"},"permissions":"all","status":"todo"}] 

#### ------------ Log Start Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTest","3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"model":{"traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"data":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"certified":true,"status":"doing","log":"","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Log Success Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6","0.89","still no error, suprah ah ah"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"model":{"traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"data":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89},"certified":true,"status":"done","log":"still no error, suprah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Query Testtuple from its key ------------
```
peer chaincode query -n mycc -c '{"Args":["query","3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"model":{"traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"data":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89},"certified":true,"status":"done","log":"still no error, suprah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Query all Testtuples ------------
```
peer chaincode query -n mycc -c '{"Args":["queryTesttuples"]}' -C myc
```
>  [{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":false,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"5608124982f25ea50a5655da1ad91eace89142e79562cbf78b6f79c5d667c9c1","log":"","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"},"permissions":"all","status":"todo"},{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6","log":"still no error, suprah ah ah","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"},"permissions":"all","status":"done"},{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"671d7c25e98b195dbe99644bcff39d9f01f0fa3b2fc09c88f8453b9e965db0f2","log":"","model":{"hash":"","storageAddress":"","traintupleKey":"f952b8514669261ca1c53ac853c3abab4a870fab395215ed8a9e7e308eb00c6b"},"permissions":"all","status":"waiting"}] 

#### ------------ Query details about a model ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModelDetails","337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"]}' -C myc
```
>  {"nonCertifiedTesttuples":[{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":false,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"5608124982f25ea50a5655da1ad91eace89142e79562cbf78b6f79c5d667c9c1","log":"","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"},"permissions":"all","status":"todo"}],"testtuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6","log":"still no error, suprah ah ah","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"},"permissions":"all","status":"done"},"traintuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"fltask":"","inModels":null,"key":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","log":"no error, ah ah ah","outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"permissions":"all","rank":0,"status":"done"}} 

#### ------------ Query all models ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModels"]}' -C myc
```
>  [{"testtuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"3a35a1c4919e8525a62885cae954690b5d2ed2e067796258a9f558064ed9d3e6","log":"still no error, suprah ah ah","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"},"permissions":"all","status":"done"},"traintuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"fltask":"","inModels":null,"key":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f","log":"no error, ah ah ah","outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"permissions":"all","rank":0,"status":"done"}},{"testtuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"671d7c25e98b195dbe99644bcff39d9f01f0fa3b2fc09c88f8453b9e965db0f2","log":"","model":{"hash":"","storageAddress":"","traintupleKey":"f952b8514669261ca1c53ac853c3abab4a870fab395215ed8a9e7e308eb00c6b"},"permissions":"all","status":"waiting"},"traintuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","data":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"fltask":"","inModels":[{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"}],"key":"f952b8514669261ca1c53ac853c3abab4a870fab395215ed8a9e7e308eb00c6b","log":"","outModel":null,"permissions":"all","rank":0,"status":"todo"}}] 

#### ------------ Query Dataset Data ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasetData","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  {"challengeKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","openerStorageAddress":"https://toto/dataset/42234/opener","owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","trainDataKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"type":"images"} 

