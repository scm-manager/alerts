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
        sh 'rm -rf content'
        sh 'mkdir content'
        // TODO
        authGit '???', 'clone https://github.com/scm-manager/website'
        sh 'mkdir -p website/content/alerts'
        script {
          image = docker.build "scmmanager/alerts:${version}"
          // TODO
          docker.withRegistry('scmmanager/alerts', 'gcloud-docker') {
            image.push()
          }
        }
      }
    }

    stage('Deployment') {
      when {
        branch 'main'
      }
      agent {
        docker {
          image 'ghcr.io/cloudogu/helm:3.4.2-1'
          args '--entrypoint=""'
          reuseNode true
        }
      }
      steps {
        // TODO
      }
    }

  }

}

def image
String version

String computeVersion() {
  def commitHashShort = sh(returnStdout: true, script: 'git rev-parse --short HEAD')
  return "${new Date().format('yyyyMMddHHmm')}-${commitHashShort}".trim()
}

void commit(String message) {
  sh "git -c user.name='Jenkins' -c user.email='jenkins@cloudogu.com' commit -m '${message}'"
}

void authGit(String credentials, String command) {
  withCredentials([
    usernamePassword(credentialsId: credentials, usernameVariable: 'AUTH_USR', passwordVariable: 'AUTH_PSW')
  ]) {
    sh "git -c credential.helper=\"!f() { echo username='\$AUTH_USR'; echo password='\$AUTH_PSW'; }; f\" ${command}"
  }
}
