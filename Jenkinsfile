pipeline {
  options {
    timestamps ()
    timeout(time: 1, unit: 'HOURS')
    buildDiscarder(logRotator(numToKeepStr: '5'))
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
          sh "go mod init chaincode"
          sh "go test chaincode/..."
        }
      }
    }

    stage('Test with substra-network') {
      steps {
        build job: 'substra-network/PR-82', parameters: [string(name: 'CHAINCODE', value: env.CHANGE_BRANCH)], propagate: true
      }
    }

  }
}
