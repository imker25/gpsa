// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

/**
* Function to set the build state on github
* 
* @param message The message to set
* @param state The state to set
 */
void setBuildStatus(String message, String state) {
  step([
      $class: "GitHubCommitStatusSetter",
      reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/imker25/gpsa"],
      contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
      errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
      statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
  ]);
}

/**
* Function to run a gradle task
*
* @param task The task to run
 */
def runBuildWorkflow(String task) {

	if (isUnix()) {
		echo "Run: \"build.sh ${task}\" on unix"
		sh "./build.sh ${task}"
	} else {
		echo "Run: \"build.bat ${task}\" on windows"
		bat ".\\build.bat ${task}"
	}
}

/** Function to clean the git repo */
def gitCleanup() {
	if (isUnix()) {
		echo "Clean workspace on unix"
		sh 'git clean -fdx'
	} else {
		echo "Clean workspace on windows"
		bat 'git clean -fx'
	}
}

/**
 Function to publish a builds output as release on github
 *
 * @param version The version string, e. g. V2.1.2-pre
 * @param text The text description of this release
 */
def publishOnGitHub(String version, String text, boolean preRelease) {
	if (isUnix()) {
		echo "${text}"
		withCredentials([string(credentialsId: 'imker25',variable: 'GITHUB_API_KEY')]) {
			sh "./build/GitHub-Release.sh ${version} \"${text}\" ${preRelease} ${GITHUB_API_KEY}"
		} 
	} else {
		throw new Exception("Can only publish on unix")
	}
}

/** The entry point of the pipeline workflow */
static void main(String[] args) {
	def labelsToRun = ["unix", "windows"]
	def buildDisplayName = ""
	def programmVersion = ""
	try {

		// Section 1: Get the build name and program version
		node("unix"){
			stage("Checkout for get build name on \"${node_name}\"") {
				echo "Checkout sources to calculate the builds name"
				checkout scm
				gitCleanup()
			}

			stage("Get build name on \"${node_name}\"") {
				echo "Get the builds name"
				runBuildWorkflow( "getBuildName")
				buildDisplayName = readFile "logs/BuildName.txt"
				echo "Set the builds display name to \"${buildDisplayName}\""
				currentBuild.displayName = 	"${buildDisplayName}"
				archiveArtifacts "logs/Version.txt"
				programmVersion = readFile "logs/Version.txt"
			}
		}

		// Section 2: Run build and test on different node types parallel (e. g. windows and unix)
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
									runBuildWorkflow( "build test")
								}


								if (isUnix()) {
									stage("Run Integration Tests on \"${node_name}\"") {
										echo "Run: build/IntegrationTests.sh"
										sh "build/IntegrationTests.sh"
									}
								}

							} finally {
								stage("Get results and artifacts") {
									junit "logs\\*.xml"
									runBuildWorkflow( "createBuildZip")
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

				parallel jobsToRun
			}
		}

		// Section 3: Publish the build output as release on github if needed
		node("unix"){
			def myBranch = "${env.BRANCH_NAME}"
			stage("Checkout for publish ${myBranch} on \"${node_name}\"") {
				echo "Checkout sources to release ${myBranch} on \"${node_name}\""
				checkout scm
				gitCleanup()
			}

			stage("Prepare release \"V${programmVersion}\" on \"${node_name}\"") {
				unarchive mapping: ['bin/' : '.']
				unarchive mapping: ['logs/' : '.']
				
			}
			if( myBranch == "master") { 
				stage("Pre release ${myBranch} on \"${node_name}\"") {
					publishOnGitHub("V${programmVersion}-pre", "Pre release of version ${programmVersion}", true)
				}
			} else if (myBranch.startsWith("release/")) {
				stage("Release ${myBranch} on \"${node_name}\"") {
					publishOnGitHub("V${programmVersion}", "Release of version ${programmVersion}", false)
				}				

			} else {
				stage("Publish on Github skipped") {
					echo "Publish on Github skipped since we running on \"${env.BRANCH_NAME}\" branch"
				}
			}
		}

		// Section 4: Set the builds status on github
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