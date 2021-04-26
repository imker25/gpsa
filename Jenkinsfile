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

static void main(String[] args) {
	def labelsToRun = ["unix", "windows"]
	def buildDisplayName = ""
	def programmVersion = ""

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
		def jobsToRun = [:]
		labelsToRun.each {  label ->
			jobsToRun["${label}"] = {
				stage("Run build and test on node with label \"${label}\"") {
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
								stage("Run Integration\"${node_name}\"") {
									sh "build/IntegrationTests.sh"
								}
							}

						} finally {
							stage("Get results and artifacts") {
								runGradle( "convertTestResults")
								junit "logs\\*.xml"
								runGradle( "createBuildZip")
								archiveArtifacts "*.zip"
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
		}
		
		parallel jobsToRun
	}

}