#!groovy
pipeline {

  options {
    buildDiscarder(logRotator(numToKeepStr: '10'))
    disableConcurrentBuilds()
  }

  agent {
    node {
      label 'docker'
    }
  }

  environment {
    HOME = "${env.WORKSPACE}"
  }

  stages {

    stage('Compute Version') {
      steps {
        script {
          version = computeVersion()
        }
      }
    }

    stage('Tests') {
      agent {
        docker {
          image 'golang:1.17.5'
          reuseNode true
        }
      }
      steps {
        sh 'go test ./...'
      }
    }

    stage('Build') {
      agent {
        docker {
          image 'golang:1.17.5'
          reuseNode true
        }
      }
      steps {
        sh 'go build -a -tags netgo -ldflags \'-w -extldflags "-static"\' -o alerts app.go'
      }
    }

    stage('Build Image') {
      when {
        branch 'main'
      }
      steps {
        dir("website") {
          git changelog: false, poll: false, branch: 'master', url: 'https://github.com/scm-manager/website'
        }
        sh 'mkdir -p website/content/alerts'
        script {
          def image = docker.build "scmmanager/alerts:${version}"
          docker.withRegistry('', 'cesmarvin-dockerhub-access-token') {
            image.push()
          }
        }
      }
    }

    stage('Deploy') {
      when {
        branch 'main'
      }
      agent {
        docker {
          image 'lachlanevenson/k8s-helm:v3.2.1'
          args  '--entrypoint=""'
        }
      }
      steps {
        withCredentials([file(credentialsId: 'helm-client-scm-manager', variable: 'KUBECONFIG')]) {
          sh "helm upgrade --install --set image.tag=${version} alerts helm/alerts"
        }
      }
    }
    
    stage('Push to GitHub') {
      when {
        branch 'main'
      }
      steps {
        authGit 'cesmarvin', "push -f --tags https://github.com/scm-manager/alerts HEAD:main" 
      }
    }

  }

  post {
    failure {
      mail to: "scm-team@cloudogu.com",
        subject: "${JOB_NAME} - Build #${BUILD_NUMBER} - ${currentBuild.currentResult}!",
        body: "Check console output at ${BUILD_URL} to view the results."
    }
  }

}

String version

String computeVersion() {
  def commitHashShort = sh(returnStdout: true, script: 'git rev-parse --short HEAD')
  return "${new Date().format('yyyyMMddHHmm')}-${commitHashShort}".trim()
}

void authGit(String credentials, String command) {
  withCredentials([
    usernamePassword(credentialsId: credentials, usernameVariable: 'AUTH_USR', passwordVariable: 'AUTH_PSW')
  ]) {
    sh "git -c credential.helper=\"!f() { echo username='\$AUTH_USR'; echo password='\$AUTH_PSW'; }; f\" ${command}"
  }
}

