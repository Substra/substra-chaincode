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
 Inputs: `Name`, `OpenerHash`, `OpenerStorageAddress`, `Type`, `DescriptionHash`, `DescriptionStorageAddress`, `ObjectiveKey`, `Permissions`
```
peer chaincode invoke -n mycc -c '{"Args":["registerDataManager",""{\"Name\":\"liver slide\",\"OpenerHash\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"OpenerStorageAddress\":\"https://toto/dataManager/42234/opener\",\"Type\":\"images\",\"DescriptionHash\":\"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee\",\"DescriptionStorageAddress\":\"https://toto/dataManager/42234/description\",\"ObjectiveKey\":\"\",\"Permissions\":\"all\"}""]}' -C myc
```
>  {"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"} 

#### ------------ Query DataManager From key ------------
Smart contract: `queryDataManager`  
 Inputs: `elementKey`
```
peer chaincode queryDataManager -n mycc -c '{"Args":["queryDataManager",""{\"Key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"}""]}' -C myc
```
>  {"objectiveKey":"","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataManager/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","opener":{"hash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/dataManager/42234/opener"},"owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","type":"images"} 

#### ------------ Add test DataSample ------------
Smart contract: `registerDataSample`  
 Inputs: `Hashes`, `DataManagerKeys`, `TestOnly`
```
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample",""{\"Hashes\":\"bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"DataManagerKeys\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"TestOnly\":\"true\"}""]}' -C myc
```
>  {"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]} 

#### ------------ Add Objective ------------
Smart contract: `registerObjective`  
 Inputs: `Name`, `DescriptionHash`, `DescriptionStorageAddress`, `MetricsName`, `MetricsHash`, `MetricsStorageAddress`, `TestDataset`, `Permissions`
```
peer chaincode invoke -n mycc -c '{"Args":["registerObjective",""{\"Name\":\"MSI classification\",\"DescriptionHash\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"DescriptionStorageAddress\":\"https://toto/objective/222/description\",\"MetricsName\":\"accuracy\",\"MetricsHash\":\"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"MetricsStorageAddress\":\"https://toto/objective/222/metrics\",\"TestDataset\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"Permissions\":\"all\"}""]}' -C myc
```
>  {"key":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"} 

#### ------------ Add Algo ------------
Smart contract: `registerAlgo`  
 Inputs: `Name`, `Hash`, `StorageAddress`, `DescriptionHash`, `DescriptionStorageAddress`, `Permissions`
```
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo",""{\"Name\":\"hog + svm\",\"Hash\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"StorageAddress\":\"https://toto/algo/222/algo\",\"DescriptionHash\":\"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca\",\"DescriptionStorageAddress\":\"https://toto/algo/222/description\",\"Permissions\":\"all\"}""]}' -C myc
```
>  {"key":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"} 

#### ------------ Add Train DataSample ------------
Smart contract: `registerDataSample`  
 Inputs: `Hashes`, `DataManagerKeys`, `TestOnly`
```
peer chaincode invoke -n mycc -c '{"Args":["registerDataSample",""{\"Hashes\":\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"DataManagerKeys\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"TestOnly\":\"false\"}""]}' -C myc
```
>  {"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]} 

#### ------------ Query DataManagers ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDataManagers"]}' -C myc
```
>  [{"objectiveKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataManager/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","opener":{"hash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/dataManager/42234/opener"},"owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","type":"images"}] 

#### ------------ Query Objectives ------------
```
peer chaincode query -n mycc -c '{"Args":["queryObjectives"]}' -C myc
```
>  [{"key":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","name":"MSI classification","description":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/description"},"metrics":{"name":"accuracy","hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"},"owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","testDataset":{"dataManagerKey":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataSampleKeys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]},"permissions":"all"}] 

#### ------------ Add Traintuple ------------
Smart contract: `createTraintuple`  
 Inputs: `AlgoKey`, `ObjectiveKey`, `InModels`, `DataManagerKey`, `DataSampleKeys`, `FLtask`, `Rank`, `Tag`
```
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple",""{\"AlgoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"ObjectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"InModels\":\"\",\"DataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"DataSampleKeys\":\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"FLtask\":\"\",\"Rank\":\"\",\"Tag\":\"\"}""]}' -C myc
```
>  {"key":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"} 

#### ------------ Add Traintuple with inModel from previous traintuple ------------
Smart contract: `createTraintuple`  
 Inputs: `AlgoKey`, `ObjectiveKey`, `InModels`, `DataManagerKey`, `DataSampleKeys`, `FLtask`, `Rank`, `Tag`
```
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple",""{\"AlgoKey\":\"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"ObjectiveKey\":\"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379\",\"InModels\":\"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687\",\"DataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"DataSampleKeys\":\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"FLtask\":\"\",\"Rank\":\"\",\"Tag\":\"\"}""]}' -C myc
```
>  {"key":"46ab1f11d49795f41e847e29e30fbd511a07dc231cf6991aa6da05cdc4a87cce"} 

#### ------------ Query Traintuples of worker with todo status ------------
```
peer chaincode invoke -n mycc -c '{"Args":["queryFilter",""{\"IndexName\":\"traintuple~worker~status\",\"Attributes\":\"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0, todo\"}""]}' -C myc
```
>  [{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"fltask":"","inModels":null,"key":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","log":"","objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"outModel":null,"permissions":"all","rank":0,"status":"todo","tag":""}] 

#### ------------ Log Start Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrain",""{\"Key\":\"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687\"}""]}' -C myc
```
>  {"algoKey":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"dataManagerKey":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataSampleKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]},"fltask":"","inModels":null,"log":"","objectiveKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","outModel":null,"perf":0,"permissions":"all","rank":0,"status":"doing","tag":""} 

#### ------------ Log Success Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain",""{\"Key\":\"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687\",\"Log\":\"no error, ah ah ah\",\"OutModel\":{\"Hash\":\"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed\",\"StorageAddress\":\"https://substrabac/model/toto\"},\"Perf\":0.9}""]}' -C myc
```
>  {"algoKey":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"dataManagerKey":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataSampleKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]},"fltask":"","inModels":null,"log":"no error, ah ah ah","objectiveKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"perf":0.9,"permissions":"all","rank":0,"status":"done","tag":""} 

#### ------------ Query Traintuple From key ------------
```
peer chaincode queryTraintuple -n mycc -c '{"Args":["queryTraintuple",""{\"Key\":\"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687\"}""]}' -C myc
```
>  {"key":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"fltask":"","inModels":null,"log":"no error, ah ah ah","objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"permissions":"all","rank":0,"status":"done","tag":""} 

#### ------------ Add Non-Certified Testtuple ------------
Smart contract: `createTesttuple`  
 Inputs: `TraintupleKey`, `DataManagerKey`, `DataSampleKeys`, `Tag`
```
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple",""{\"TraintupleKey\":\"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687\",\"DataManagerKey\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"DataSampleKeys\":\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"Tag\":\"\"}""]}' -C myc
```
>  {"key":"cee7d90187be57dae0f83195302abd4c446c5e52fb49abb6b3e1607b970b2577"} 

#### ------------ Add Certified Testtuple ------------
Smart contract: `createTesttuple`  
 Inputs: `TraintupleKey`, `DataManagerKey`, `DataSampleKeys`, `Tag`
```
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple",""{\"TraintupleKey\":\"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687\",\"DataManagerKey\":\"\",\"DataSampleKeys\":\"\",\"Tag\":\"\"}""]}' -C myc
```
>  {"key":"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53"} 

#### ------------ Add Testtuple with not done traintuple ------------
Smart contract: `createTesttuple`  
 Inputs: `TraintupleKey`, `DataManagerKey`, `DataSampleKeys`, `Tag`
```
peer chaincode invoke -n mycc -c '{"Args":["createTesttuple",""{\"TraintupleKey\":\"46ab1f11d49795f41e847e29e30fbd511a07dc231cf6991aa6da05cdc4a87cce\",\"DataManagerKey\":\"\",\"DataSampleKeys\":\"\",\"Tag\":\"\"}""]}' -C myc
```
>  {"key":"88d515424fddbc49e8f0202cecaff056f77a49faa59688346eab2e68c9dd8c46"} 

#### ------------ Query Testtuples of worker with todo status ------------
```
peer chaincode invoke -n mycc -c '{"Args":["queryFilter",""{\"IndexName\":\"testtuple~worker~status\",\"Attributes\":\"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0, todo\"}""]}' -C myc
```
>  [{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53","log":"","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"todo","tag":""},{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":false,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"cee7d90187be57dae0f83195302abd4c446c5e52fb49abb6b3e1607b970b2577","log":"","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"todo","tag":""}] 

#### ------------ Log Start Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTest",""{\"Key\":\"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53\"}""]}' -C myc
```
>  {"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"log":"","model":{"traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"doing","tag":""} 

#### ------------ Log Success Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest",""{\"Key\":\"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53\",\"Log\":\"no error, ah ah ah\",\"Perf\":0.9}""]}' -C myc
```
>  {"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"log":"no error, ah ah ah","model":{"traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"done","tag":""} 

#### ------------ Query Testtuple from its key ------------
```
peer chaincode queryTesttuple -n mycc -c '{"Args":["queryTesttuple",""{\"Key\":\"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53\"}""]}' -C myc
```
>  {"key":"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53","algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"log":"no error, ah ah ah","model":{"traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"done","tag":""} 

#### ------------ Query all Testtuples ------------
```
peer chaincode queryTesttuples -n mycc -c '{"Args":["queryTesttuples"]}' -C myc
```
>  [{"key":"88d515424fddbc49e8f0202cecaff056f77a49faa59688346eab2e68c9dd8c46","algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"log":"","model":{"traintupleKey":"46ab1f11d49795f41e847e29e30fbd511a07dc231cf6991aa6da05cdc4a87cce","hash":"","storageAddress":""},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"waiting","tag":""},{"key":"cee7d90187be57dae0f83195302abd4c446c5e52fb49abb6b3e1607b970b2577","algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"certified":false,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"log":"","model":{"traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"todo","tag":""},{"key":"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53","algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"log":"no error, ah ah ah","model":{"traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"done","tag":""}] 

#### ------------ Query details about a model ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModelDetails",""{\"Key\":\"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687\"}""]}' -C myc
```
>  {"nonCertifiedTesttuples":[{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":false,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"cee7d90187be57dae0f83195302abd4c446c5e52fb49abb6b3e1607b970b2577","log":"","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"todo","tag":""}],"testtuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53","log":"no error, ah ah ah","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"done","tag":""},"traintuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"fltask":"","inModels":null,"key":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","log":"no error, ah ah ah","objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"permissions":"all","rank":0,"status":"done","tag":""}} 

#### ------------ Query all models ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModels"]}' -C myc
```
>  [{"testtuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"88d515424fddbc49e8f0202cecaff056f77a49faa59688346eab2e68c9dd8c46","log":"","model":{"hash":"","storageAddress":"","traintupleKey":"46ab1f11d49795f41e847e29e30fbd511a07dc231cf6991aa6da05cdc4a87cce"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"waiting","tag":""},"traintuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"fltask":"","inModels":[{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"}],"key":"46ab1f11d49795f41e847e29e30fbd511a07dc231cf6991aa6da05cdc4a87cce","log":"","objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"outModel":null,"permissions":"all","rank":0,"status":"todo","tag":""}},{"testtuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"certified":true,"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"key":"6581232b890c5f80522ca480b815d3340dffde8924d863bfd1dfedb1841a2d53","log":"no error, ah ah ah","model":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto","traintupleKey":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"},"objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"permissions":"all","status":"done","tag":""},"traintuple":{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","dataset":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"fltask":"","inModels":null,"key":"8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687","log":"no error, ah ah ah","objective":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","storageAddress":"https://toto/objective/222/metrics"}},"outModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"permissions":"all","rank":0,"status":"done","tag":""}}] 

#### ------------ Query Dataset ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDataset",""{\"Key\":\"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"}""]}' -C myc
```
>  {"objectiveKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataManager/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","opener":{"hash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/dataManager/42234/opener"},"owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","type":"images","trainDataSampleKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"testDataSampleKeys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]} 

#### ------------ Update Data Sample with new data manager ------------
```
peer chaincode invoke -n mycc -c '{"Args":["updateDataSample",""{\"Hashes\":\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\",\"DataManagerKeys\":\"38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee\"}""]}' -C myc
```
>  {"key":"{\"keys\": [\"aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc\"]}"} 

#### ------------ Query the new Dataset ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDataset",""{\"Key\":\"38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee\"}""]}' -C myc
```
>  {"objectiveKey":"","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataManager/42234/description"},"key":"38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee","name":"liver slide","opener":{"hash":"38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee","storageAddress":"https://toto/dataManager/42234/opener"},"owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","type":"images","trainDataSampleKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"testDataSampleKeys":[]} 

