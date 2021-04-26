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
		sh "gradle ${task}"
	} else {
		bat "gradle ${task}"
	}
}

static void main(String[] args) {

	node("unix"){
		stage("Checkout for get build name on \"${node_name}\"") {
			echo "Checkout sources to calculate the builds name"
			checkout scm
			sh 'git clean -fdx'
		}

		stage("Get build name on \"${node_name}\"") {
			runGradle( "getBuildName")
			currentBuild.displayName = readFile "logs/BuildName.txt"		
		}
	}

}