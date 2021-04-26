// Copyright 2019 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.
def programmVersion setBuildStatus(String message, String state) {
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

static void main(String[] args) {
	def buildDisplayName = ""
	node("unix"){
		stage("Checkout for get build name on \"${node_name}\"") {
			echo "Checkout sources to calculate the builds name"
			checkout scm
			sh 'git clean -fdx'
		}

		stage("Get build name on \"${node_name}\"") {
			echo "Get the builds name"
			runGradle( "getBuildName")
			buildDisplayName = readFile "logs/BuildName.txt"
			echo "Set the builds display name to \"${buildDisplayName}\""
			currentBuild.displayName = 	"${buildDisplayName}"
		}
	}

}