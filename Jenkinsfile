#!groovy

import groovy.transform.Field

@Field String email_to = 'sw@platinasystems.com'
@Field String email_from = 'jenkins-bot@platinasystems.com'
@Field String email_reply_to = 'no-reply@platinasystems.com'

pipeline {
    agent any
    stages {
	stage('Build') {
	    steps {
		echo "Running go vet on goes-bmc..."
		sh 'set +x; go vet && go vet ./cmd/...'
		echo "Running go test on goes-bmc..."
		sh 'set +x; go test && go test ./cmd/...'
		echo "Building goes-bmc"
		sh 'set +x; go build'
	    }
	}
    }

    post {
	success {
	    mail body: "GOES-BMC build ok: ${env.BUILD_URL}",
		from: email_from,
		replyTo: email_reply_to,
		subject: 'GOES-BMC build ok',
		to: email_to
	}
	failure {
	    cleanWs()
	    mail body: "GOES-BMC build error: ${env.BUILD_URL}",
		from: email_from,
		replyTo: email_reply_to,
		subject: 'GOES-BMC BUILD FAILED',
		to: email_to
	}
    }
}
