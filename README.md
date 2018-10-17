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
- `registerData`
- `registerAlgo`
- `createTraintuple`
- `logStartTrainTest`
- `logSuccessTrain`
- `logSuccessTest`
- `logFailTrainTest`
- `query`
- `queryProblems`
- `queryAlgo`
- `queryTraintuples`
- `queryDatasets`
- `queryModel`
- `queryModelTraintuples`
- `queryDatasetData`
- `queryFilter`

### Examples 
#### ------------ Add Dataset ------------
Smart contract: `registerDataset`
Inputs: `Name`, `OpenerHash`, `OpenerStorageAddress`, `Type`, `DescriptionHash`, `DescriptionStorageAddress`, `ChallengeKeys`, `Permissions`
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
>  {"name":"liver slide","openerStorageAddress":"https://toto/dataset/42234/opener","size":0,"nbData":0,"type":"images","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","challengeKeys":[],"permissions":"all"} 

#### ------------ Add test Data ------------
Smart contract: `registerData`
Inputs: `Hashes`, `DatasetKey`, `Size`, `TestOnly`
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","100","true"]}' -C myc
```
>  bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc 

#### ------------ Add Challenge ------------
Smart contract: `registerChallenge`
Inputs: `Name`, `DescriptionHash`, `DescriptionStorageAddress`, `MetricsName`, `MetricsHash`, `MetricsStorageAddress`, `TestDataKeys`, `Permissions`
```
peer chaincode invoke -n mycc -c '{"Args":["registerChallenge","MSI classification","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","https://toto/challenge/222/description","accuracy","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","https://toto/challenge/222/metrics","bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","all"]}' -C myc
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
Inputs: `Hashes`, `DatasetKey`, `Size`, `TestOnly`
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","100","false"]}' -C myc
```
>  aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc 

#### ------------ Query Datasets ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasets"]}' -C myc
```
>  [{"challengeKeys":["5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"],"description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","nbData":4,"openerStorageAddress":"https://toto/dataset/42234/opener","owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","size":200,"type":"images"}] 

#### ------------ Query Challenges ------------
```
peer chaincode query -n mycc -c '{"Args":["queryChallenges"]}' -C myc
```
>  [{"descriptionStorageAddress":"https://toto/challenge/222/description","key":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","name":"accuracy","storageAddress":"https://toto/challenge/222/metrics"},"name":"MSI classification","owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","testDataKeys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}] 

#### ------------ Add Traintuple ------------
Smart contract: `createTraintuple`
Inputs: `AlgoKey`, `StartModelKey`, `TrainDataKeys`
```
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d 

#### ------------ Query Traintuples of worker with todo status ------------
```
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","traintuple~trainWorker~status","bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0, todo"]}' -C myc
```
>  [{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","endModel":null,"key":"655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d","log":"","permissions":"all","startModel":null,"status":"todo","testData":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"trainData":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"}}] 

#### ------------ Log Start Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrainTest","655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d","training"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":null,"endModel":null,"trainData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"testData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"status":"training","log":"","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Log Success Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d","eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed, https://substrabac/model/toto","0.9","no error, ah ah ah"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":null,"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"testData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"status":"trained","log":"no error, ah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Log Start Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrainTest","655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d","testing"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":null,"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"testData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0},"status":"testing","log":"no error, ah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Log Success Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d","0.89","still no error, suprah ah ah"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":null,"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"testData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89},"status":"done","log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Query Traintuple From key ------------
```
peer chaincode query -n mycc -c '{"Args":["query","655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":null,"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"testData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89},"status":"done","log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Query Model ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModel","eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"name":"hog + svm","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":null,"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9},"testData":{"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89},"status":"done","log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"} 

#### ------------ Query Model Traintuples ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModelTraintuples","eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"]}' -C myc
```
>  {"algo":{"challengeKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","description":{"hash":"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca","storageAddress":"https://toto/algo/222/description"},"key":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","storageAddress":"https://toto/algo/222/algo"},"traintuples":[{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","storageAddress":"https://toto/algo/222/algo"},"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"key":"655f409da18f1f79dcdef637c491146528958d6a3842e111b0797ecbfa1d2c8d","log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","startModel":null,"status":"done","testData":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.89,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"},"trainData":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":0.9,"worker":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"}}]} 

#### ------------ Query Dataset Data ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasetData","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  {"challengeKeys":["5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"],"description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","nbData":4,"openerStorageAddress":"https://toto/dataset/42234/opener","owner":"bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0","permissions":"all","size":200,"trainDataKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"type":"images"} 

