// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
def programmVersion
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
    agent {
		label "awaiter"
	}
	options { skipDefaultCheckout() }
	environment {
        GITHUB_API_KEY = credentials('imker25')
        
    }
	
    stages {
		stage('Get Build Name') {
			steps ('Calculate the name') {
				checkout scm
				sh 'gradle getBuildName'
				script {
					currentBuild.displayName = readFile "logs/BuildName.txt"
				}
			}
		}
		

        stage('Build and test the gpsa project') {
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

			// ToDo:
	// Write new stages that get the artifacts back
	// https://jenkins.io/doc/pipeline/steps/workflow-basic-steps/#-unarchive-copy-archived-artifacts-into-the-workspace
	// and uploads them
	//
	// use the when expression
	// https://jenkins.io/doc/book/pipeline/syntax/#when
	// to figure out what branch
	// do a pre release on master
	// do a release on feture braches


		stage('Publish master') {
			when {
                branch 'master'
            }
			steps ('Do a pre release') { 
				unarchive mapping: ['bin/' : '.']
				script {
					programmVersion = readFile "logs/Version.txt"
				}
				
				sh "./build/GitHub-Release.sh V${programmVersion}-pre \"Pre release of version ${programmVersion}\" true ${GITHUB_API_KEY}"
				
			}
		}

		stage('Publish release') {
			when {
                branch 'release/**'
            }
			steps ('Do a release') { 
				unarchive mapping: ['bin/' : '.']
				script {
					programmVersion = readFile "logs/Version.txt"
				}
				
				sh "./build/GitHub-Release.sh V${programmVersion} \"Release of version ${programmVersion}\" false ${GITHUB_API_KEY}"
				
			}
		}
    }
	post ('Publish build result on GitHub') {

		always {
			archiveArtifacts "logs/*.json"
		}

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
