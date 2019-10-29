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
							bat 'gradle createBuildZip'
							archiveArtifacts "*.zip"
						}
					}	
					
				}
                
        
                stage('Run on Linux') {
                    agent {
                        label "unix"
                    }
					stages{
						stage('Prepare linux worspace'){
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
							sh 'gradle createBuildZip'	
							archiveArtifacts "*.zip"													
						}						
					}
                }
            }
        }
    }
}