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
- `queryModels`
- `queryDatasets`
- `queryModel`
- `queryModelTraintuples`
- `queryDatasetData`

### Examples 

#### ------------ Add Dataset ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerDataset","liver slide","do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","https://toto/dataset/42234/opener","images","8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","https://toto/dataset/42234/description","","all"]}' -C myc
```
>  {"name":"liver slide","size":0,"nbData":0,"type":"images","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","challengeKeys":null,"permissions":"all"}
#### ------------ Add test Data ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, da2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataset_do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","100","true"]}' -C myc
```
>  data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, data_da2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, 
#### ------------ Add Challenge ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerChallenge","MSI classification","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","https://toto/challenge/222/description","accuracy","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","https://toto/challenge/222/metrics","data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","all"]}' -C myc
```
>  {"name":"MSI classification","descriptionStorageAddress":"https://toto/challenge/222/description","metrics":{"name":"accuracy","hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"},"owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","testDataKeys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"permissions":"all"}
#### ------------ Add Algo ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerAlgo","hog + svm","fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","https://toto/algo/222/algo","e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca","https://toto/algo/222/description","challenge_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","all"]}' -C myc
```
>  {"name":"hog + svm","storageAddress":"https://toto/algo/222/algo","description":{"hash":"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca","storageAddress":"https://toto/algo/222/description"},"owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","challengeKey":"challenge_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","permissions":""}
#### ------------ Add Train Data ------------
```
peer chaincode invoke -n mycc -c '{"Args":["registerData","aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","dataset_do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","100","false"]}' -C myc
```
#### ------------ Query Datasets ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasets"]}' -C myc
```
>  [{"challengeKeys":["challenge_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"],"description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"key":"dataset_do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","name":"liver slide","nbData":4,"owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","permissions":"all","size":200,"type":"images"}]
#### ------------ Add Traintuple ------------
```
peer chaincode invoke -n mycc -c '{"Args":["createTraintuple","challenge_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","algo_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","algo_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc, data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  traintuple_92625570a513f7c34cacbf4d2383972f07e1118d9a2b8c33c9398ff29e892a4b
#### ------------ Log Start Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrainTest","traintuple_92625570a513f7c34cacbf4d2383972f07e1118d9a2b8c33c9398ff29e892a4b","training"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":null,"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"status":"training","rank":0,"perf":0,"log":"","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
#### ------------ Log Success Training ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTrain","traintuple_92625570a513f7c34cacbf4d2383972f07e1118d9a2b8c33c9398ff29e892a4b","modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod, https://substrabac/model/toto","data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.90, data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.91","no error, ah ah ah"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"status":"trained","rank":0,"perf":0,"log":"no error, ah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
#### ------------ Log Start Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logStartTrainTest","traintuple_92625570a513f7c34cacbf4d2383972f07e1118d9a2b8c33c9398ff29e892a4b","testing"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":null},"status":"testing","rank":0,"perf":0,"log":"no error, ah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
#### ------------ Log Success Testing ------------
```
peer chaincode invoke -n mycc -c '{"Args":["logSuccessTest","traintuple_92625570a513f7c34cacbf4d2383972f07e1118d9a2b8c33c9398ff29e892a4b","data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.90, data_da2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc:0.91","0.99","still no error, suprah ah ah"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9]},"status":"testing","rank":0,"perf":0.99,"log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
#### ------------ Query Traintuple From key ------------
```
peer chaincode query -n mycc -c '{"Args":["query","traintuple_92625570a513f7c34cacbf4d2383972f07e1118d9a2b8c33c9398ff29e892a4b"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9]},"status":"testing","rank":0,"perf":0.99,"log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}

#### ------------ Query Model ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModel","modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod"]}' -C myc
```
>  {"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9]},"status":"testing","rank":0,"perf":0.99,"log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
#### ------------ Query Model Traintuples ------------
```
peer chaincode query -n mycc -c '{"Args":["queryModelTraintuples","modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod"]}' -C myc
```
>  {"algo_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc":{"name":"hog + svm","storageAddress":"https://toto/algo/222/algo","description":{"hash":"e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca","storageAddress":"https://toto/algo/222/description"},"owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","challengeKey":"challenge_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","permissions":""},"traintuple_92625570a513f7c34cacbf4d2383972f07e1118d9a2b8c33c9398ff29e892a4b":{"challenge":{"hash":"5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379","metrics":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","storageAddress":"https://toto/challenge/222/metrics"}},"algo":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"startModel":{"hash":"fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","storageAddress":"https://toto/algo/222/algo"},"endModel":{"hash":"modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod","storageAddress":"https://substrabac/model/toto"},"trainData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9,0.91]},"testData":{"worker":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","keys":["data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"],"openerHash":"do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","perf":[0.9]},"status":"testing","rank":0,"perf":0.99,"log":"no error, ah ah ahstill no error, suprah ah ah","permissions":"all","creator":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}}
#### ------------ Query Dataset Data ------------
```
peer chaincode query -n mycc -c '{"Args":["queryDatasetData","dataset_do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}' -C myc
```
>  {"dataset_do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc":{"name":"liver slide","size":200,"nbData":4,"type":"images","description":{"hash":"8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee","storageAddress":"https://toto/dataset/42234/description"},"owner":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","challengeKeys":["challenge_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"],"permissions":"all"},"trainDataKeys":["data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc","data_aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"]}

