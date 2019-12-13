# substra-chaincode

[![Build Status](https://travis-ci.org/SubstraFoundation/substra-chaincode.svg?branch=master)](https://travis-ci.org/SubstraFoundation/substra-chaincode)

Chaincode for the Substra platform

> :warning: This project is under active development. Please wait a while before using it...

## License

This project is developed under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](./LICENSE) file.

## Local development

### Prerequisites

go version 1.12.x

### Run the tests

```
cd chaincode
go test
```

## Devmode

See [chaincode-docker-devmode](./chaincode-docker-devmode/README.rst)

## Documentation


Note for internal use only: See the [technical specifications](https://github.com/SubstraFoundation/substra-spec/blob/master/technical_spec_substra.md#smartcontract).

### Implemented smart contracts

- `createAggregatetuple`
- `createCompositeTraintuple`
- `createComputePlan`
- `createTesttuple`
- `createTraintuple`
- `logFailAggregate`
- `logFailCompositeTrain`
- `logFailTest`
- `logFailTrain`
- `logStartAggregate`
- `logStartCompositeTrain`
- `logStartTest`
- `logStartTrain`
- `logSuccessAggregate`
- `logSuccessCompositeTrain`
- `logSuccessTest`
- `logSuccessTrain`
- `queryAggregateAlgo`
- `queryAggregateAlgos`
- `queryAggregatetuple`
- `queryAggregatetuples`
- `queryAlgo`
- `queryAlgos`
- `queryCompositeAlgo`
- `queryCompositeAlgos`
- `queryCompositeTraintuple`
- `queryCompositeTraintuples`
- `queryComputePlan`
- `queryComputePlans`
- `queryDataManager`
- `queryDataManagers`
- `queryDataSamples`
- `queryDataset`
- `queryFilter`
- `queryModelDetails`
- `queryModels`
- `queryNodes`
- `queryObjective`
- `queryObjectiveLeaderboard`
- `queryObjectives`
- `queryTesttuple`
- `queryTesttuples`
- `queryTraintuple`
- `queryTraintuples`
- `registerAggregateAlgo`
- `registerAlgo`
- `registerCompositeAlgo`
- `registerDataManager`
- `registerDataSample`
- `registerNode`
- `registerObjective`
- `updateDataManager`
- `updateDataSample`

### Examples

See the [full list of examples](./EXAMPLES.md)
