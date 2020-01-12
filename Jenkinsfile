// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

void setBuildStatus(String message, String state) {
  step([
      $class: "GitHubCommitStatusSetter",
      reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/imker25/gpsa"],
      contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
      errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
      statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
  ]);
}

pipeline {
    agent none
	options { skipDefaultCheckout() }

    stages {
        stage('Build, test and deploy the gpsa project') {
            parallel {
                stage('Run on Windows') {
                    agent {
                        label "windows"
                    }
					stages{
						stage('Prepare windows workspace'){
							steps ('Checkout') {
								checkout scm
							}
						}
						stage('Create windows binaries') {

							steps ('Build') {
								bat 'gradle build'
							}
						}
						stage('Test windows binaries') {
							steps('Run') {
								bat 'gradle test'
							}

						}
					}
					post('Deploy windows results') {
        				always {
							bat 'gradle convertTestResults'
							junit "logs\\*.xml"
							bat 'gradle createBuildZip'
							archiveArtifacts "*.zip"
							archiveArtifacts "bin\\gpsa.exe"
						}
					}

				}


                stage('Run on Linux') {
                    agent {
                        label "unix"
                    }
					stages{
						stage('Prepare linux workspace'){
							steps ('Checkout') {
								checkout scm
							}
						}
						stage('Create linux binaries') {
							steps ('Build') {
								sh 'gradle build'
							}
						}
						stage('Test linux binaries') {
							steps('Run') {
								sh 'gradle test'
							}

						}
					}
					post('Deploy linux results') {
        				always {
							sh 'gradle convertTestResults'
							junit "logs/*.xml"
							sh 'gradle createBuildZip'
							archiveArtifacts "*.zip"
							archiveArtifacts "bin/gpsa"
						}
					}
                }
            }
        }
    }
	post ('Publish build result on GitHub') {
		success {
			setBuildStatus("Build complete", "SUCCESS");
		}

		failure {
			setBuildStatus("Build complete", "FAILURE");
		}

		unstable {
			setBuildStatus("Build complete", "UNSTABLE");
		}
	}
}
