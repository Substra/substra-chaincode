pipeline {
  options {
    timestamps ()
    timeout(time: 1, unit: 'HOURS')
    buildDiscarder(logRotator(numToKeepStr: '5'))
  }

  parameters {
    booleanParam(name: 'E2E', defaultValue: false, description: 'Launch E2E test')
    string(name: 'BACKEND', defaultValue: 'dev', description: 'substrabac branch')
    string(name: 'CLI', defaultValue: 'dev', description: 'substra-cli branch')
  }

  agent none

  stages {
    stage('Abort previous builds'){
      steps {
        milestone(Integer.parseInt(env.BUILD_ID)-1)
        milestone(Integer.parseInt(env.BUILD_ID))
      }
    }

    stage('Test') {
      agent {
        kubernetes {
          label 'substra-chaincode'
          defaultContainer 'go'
          yaml """
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: go
                image: golang:1.12
                command: [cat]
                tty: true
            """
        }
      }

      steps {
        dir("chaincode") {
          sh "go test chaincode/..."
        }
      }
    }

    stage('Test with substra-network') {
      when {
        expression { return params.E2E }
      }

      steps {
        build job: 'substra-network/dev', parameters: [string(name: 'CHAINCODE', value: env.CHANGE_BRANCH),
                                                       string(name: 'BACKEND', value: params.BACKEND),
                                                       string(name: 'CLI', value: params.CLI)], propagate: true
      }
    }

  }
}
