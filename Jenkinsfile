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

def runGradle(String task) {

	if (isUnix()) {
		echo "Run: \"gradle ${task}\" on unix"
		sh "gradle ${task}"
	} else {
		echo "Run: \"gradle ${task}\" on windows"
		bat "gradle ${task}"
	}
}

def gitCleanup() {
		if (isUnix()) {
		echo "Clean workspace on unix"
		sh 'git clean -fdx'
	} else {
		echo "Clean workspace on windows"
		bat 'git clean -fdx'
	}
}

def publishOnGitHub(String version, String text) {
	if (isUnix()) {
		echo "${text}"
		withCredentials([usernameColonPassword(credentialsId: 'imker25',variable: 'GITHUB_API_KEY')]) {
			sh "./build/GitHub-Release.sh V${version}-pre \"${text}\" true ${GITHUB_API_KEY}"
		}
	} else {
		throw new Exception("Can only publish on unix")
	}
}

static void main(String[] args) {
	def labelsToRun = ["unix", "windows"]
	def buildDisplayName = ""
	def programmVersion = ""
	try {
		node("unix"){
			stage("Checkout for get build name on \"${node_name}\"") {
				echo "Checkout sources to calculate the builds name"
				checkout scm
				gitCleanup()
			}

			stage("Get build name on \"${node_name}\"") {
				echo "Get the builds name"
				runGradle( "getBuildName")
				buildDisplayName = readFile "logs/BuildName.txt"
				echo "Set the builds display name to \"${buildDisplayName}\""
				currentBuild.displayName = 	"${buildDisplayName}"
			}
		}

		node("awaiter") {
			stage("Run build and test on nodes with labels ${labelsToRun}") {
				def jobsToRun = [:]
				labelsToRun.each {  label ->
					jobsToRun["${label}"] = {
					
						node("${label}"){
							try {
								stage("Checkout for build and test on \"${node_name}\"") {
									echo "Checkout sources on \"${node_name}\""
									checkout scm
									gitCleanup()
								}

								stage("Build on \"${node_name}\"") {
									runGradle( "build")
								}

								stage("Test on \"${node_name}\"") {
									runGradle( "test")
								}

								if (isUnix()) {
									stage("Run Integration Tests on \"${node_name}\"") {
										echo "Run: build/IntegrationTests.sh"
										sh "build/IntegrationTests.sh"
									}
								}

							} finally {
								stage("Get results and artifacts") {
									runGradle( "convertTestResults")
									junit "logs\\*.xml"
									runGradle( "createBuildZip")
									archiveArtifacts "*.zip"
									archiveArtifacts "logs/Version.txt"
									if (isUnix()) {
										archiveArtifacts "bin/gpsa"
									} else {
										archiveArtifacts "bin/gpsa.exe"
									}
								}
							}
						}
					}
				}

				parallel jobsToRun
			}
		}

		node("unix"){
			def myBranch = "${env.BRANCH_NAME}"
			stage("Checkout for publish ${myBranch} on \"${node_name}\"") {
				echo "Checkout sources to release ${myBranch} on \"${node_name}\""
				checkout scm
				gitCleanup()
			}

			stage("Prepare release on \"${node_name}\"") {
				unarchive mapping: ['bin/' : '.']
				unarchive mapping: ['logs/' : '.']
				programmVersion = readFile "logs/Version.txt"
			}
			if( myBranch == "master") {
				stage("Pre release ${myBranch} on \"${node_name}\"") {
					publishOnGitHub("${programmVersion}", "Pre release of version ${programmVersion}")
				}
			} else if (myBranch.startsWith("release/")) {
				stage("Release ${myBranch} on \"${node_name}\"") {
					publishOnGitHub("${programmVersion}", "Release of version ${programmVersion}")
				}				

			} else {
				stage("Publish on Github skipped") {
					echo "Publish on Github skipped since we running on \"${env.BRANCH_NAME}\" branch"
				}
			}
		}

		node("awaiter") {
			stage("Finanlize build") {
				setBuildStatus("Build complete", "SUCCESS")
			}
		}
	} catch(error) {
		node("awaiter") {
			stage("Finanlize build wiht FAILURE") {
				setBuildStatus("Build complete", "FAILURE")
			}
			throw error
		}
	}

}