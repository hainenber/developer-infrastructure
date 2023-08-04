pipelineJob('hetman-ci') {
    properties {
        pipelineTriggers {
            triggers {
                cron {
                    spec('H/5 * * * *')
                }
            }
        }
    }

    definition {
        cps {
            script("""
                podTemplate(yaml: '''
                    apiVersion: v1
                    kind: Pod
                    spec:
                      containers:
                        - name: golang
                          image: golang:1.20
                          command: [ 'sleep' ]
                          args: [ '99d' ]
                        - name: kaniko
                          image: gcr.io/kaniko-project/executor:debug
                          command: [ 'sleep' ]
                          args: [ '99999' ]
                '''
                ) {
                    node(POD_LABEL) {
                        // Checkout
                        git(
                            branch: 'main',
                            credentialsId: 'hainenber_personal_access_token',
                            url: 'https://github.com/hainenber/hetman.git'
                        )

                        // Test
                        container('golang') {
                            stage('Test') {
                                sh 'go test -timeout 2m -cover ./...'
                            }
                        }

                        // Build image
                        container('kaniko') {
                            stage('Build image') {
                                // Make the output directory.
                                sh "mkdir -p output"

                                // Workaround to tell Kaniko keeping files between multi-stage builds
                                sh 'touch /kaniko/.docker/config.json'

                                // Build and save built image to tar archive
                                withEnv(["DOCKER_CONFIG=/kaniko/.docker/"]) {
                                    sh '/kaniko/executor --context `pwd` --no-push --tar-path=./output/hetman.tar'
                                }
                            }
                        }

                        // container('syft') {
                        //     stage('SBOM generation') {
                        //         sh 'syft hetman.tar -o json=./output/sbom.syft.json -o table' 
                        //     }
                        // }

                        
                        // Vulnerability scanning with anchore/grype
                        sh 'grype sbom:./output/sbom.syft.json -o json=./output/vuln.grype.json -o table'

                        // Archive validation results as artifacts
                        archiveArtifacts artifacts: 'output/*.json'

                        // Push to private Registry
                        docker.withRegistry('http://docker-registry:32001') {
                            hetmanImage.push()
                        }
                    }
                }
            """)
        }
    }
}
