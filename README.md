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
```
peer chaincode invoke -n mycc -c '{"Args":["registerDataset","liver slide","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","https://toto/dataset/42234/opener","images","8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","https://toto/dataset/42234/description","","all"]}' -C myc
```
>  da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc 

#### ------------ Query Dataset From key ------------
```
peer chaincode query -n mycc -c '{"Args":["query","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  {"name":"liver slide","openerStorageAddress":"https://toto/dataset/42234/opener","size":0,"nbData":0,"type":"images","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","challengeKeys":null,"permissions":"all"} 

#### ------------ Add test Data ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","100","true"]}' -C myc
```
>  bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc,  

#### ------------ Add Challenge ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerChallenge","MSI classification","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","https://toto/challenge/222/description","accuracy","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","https://toto/challenge/222/metrics","bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","all"]}' -C myc
```
>  da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc 

#### ------------ Add Algo ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo","hog + svm","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","https://toto/algo/222/algo","e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca","https://toto/algo/222/description","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","all"]}' -C myc
```
>  fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc 

#### ------------ Add Train Data ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","100","false"]}' -C myc
```
>  aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc,  

#### ------------ Query Datasets ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasets"]}' -C myc
```
>  [{"challengeKeys":["5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"],"description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","nbData":4,"openerStorageAddress":"https://toto/dataset/42234/opener","owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","permissions":"all","size":200,"type":"images"}] 

#### ------------ Query Challenges ------------
```
peer chaincode query -n mycc -c '{"Args":["queryChallenges"]}' -C myc
```
>  [{"descriptionStorageAddress":"https://toto/challenge/222/description","key":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","name":"accuracy","storageAddress":"https://toto/challenge/222/metrics"},"name":"MSI classification","owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","permissions":"all","testDataKeys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}] 

#### ------------ Add Traintuple ------------
```
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  e1d91c9dd5c15f99c18f743d9e4b9bba3b5f5a7370977ac877f1214fe87bfca7 

#### ------------ Query Traintuples of worker with todo status ------------
```
peer chaincode invoke -n mycc -c '{"Args":["queryFilter","traintuple~trainWorker~status","e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855, todo"]}' -C myc
```
>  [{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","endModel":null,"log":"","perf":0,"permissions":"all","rank":0,"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"status":"todo","testData":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null,"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},"trainData":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null,"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}}] 

#### ------------ Log Start Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrainTest","e1d91c9dd5c15f99c18f743d9e4b9bba3b5f5a7370977ac877f1214fe87bfca7","training"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":null,"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"status":"training","rank":0,"perf":0,"log":"","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"} 

#### ------------ Log Success Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","e1d91c9dd5c15f99c18f743d9e4b9bba3b5f5a7370977ac877f1214fe87bfca7","eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed, https://substrabac/model/toto","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.90, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.91","no error, ah ah ah"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"status":"trained","rank":0,"perf":0,"log":"no error, ah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"} 

#### ------------ Log Start Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrainTest","e1d91c9dd5c15f99c18f743d9e4b9bba3b5f5a7370977ac877f1214fe87bfca7","testing"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"status":"testing","rank":0,"perf":0,"log":"no error, ah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"} 

#### ------------ Log Success Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","e1d91c9dd5c15f99c18f743d9e4b9bba3b5f5a7370977ac877f1214fe87bfca7","bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.90, bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.91","0.99","still no error, suprah ah ah"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"status":"done","rank":0,"perf":0.99,"log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"} 

#### ------------ Query Traintuple From key ------------
```
peer chaincode query -n mycc -c '{"Args":["query","e1d91c9dd5c15f99c18f743d9e4b9bba3b5f5a7370977ac877f1214fe87bfca7"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"status":"done","rank":0,"perf":0.99,"log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"} 

#### ------------ Query Model ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModel","eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"status":"done","rank":0,"perf":0.99,"log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"} 

#### ------------ Query Model Traintuples ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModelTraintuples","eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"]}' -C myc
```
>  {"algo":{"challengeKey":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","description":{"hash":"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca","storageAddress":"https://toto/algo/222/description"},"key":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"hog + svm","owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","permissions":"all","storageAddress":"https://toto/algo/222/algo"},"traintuples":[{"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","endModel":{"hash":"eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed","storageAddress":"https://substrabac/model/toto"},"key":"e1d91c9dd5c15f99c18f743d9e4b9bba3b5f5a7370977ac877f1214fe87bfca7","log":"no error, ah ah ahstill no error, suprah ah ah","perf":0.99,"permissions":"all","rank":0,"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"status":"done","testData":{"keys":["bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91],"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},"trainData":{"keys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91],"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}}]} 

#### ------------ Query Dataset Data ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasetData","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  {"dataset":{"challengeKeys":["5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"],"description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"key":"da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","nbData":4,"openerStorageAddress":"https://toto/dataset/42234/opener","owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","permissions":"all","size":200,"type":"images"},"trainDataKeys":["aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]} 

[31m2018-08-29 15:14:58.304 CEST [mock] HasNext -> ERRO 001[0m HasNext() couldn't get Current
[31m2018-08-29 15:14:58.304 CEST [mock] HasNext -> ERRO 002[0m HasNext() couldn't get Current
PASS
ok  	github.com/SubstraFoundation/substra-chaincode/chaincode	0.027s
