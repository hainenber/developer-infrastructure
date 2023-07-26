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
                node('docker') {
                    // Checkout
                    git(
                        branch: 'main',
                        credentialsId: 'hainenber_personal_access_token',
                        url: 'https://github.com/hainenber/hetman.git'
                    )

                    // Test
                    sh 'go test -timeout 2m -cover ./...'

                    // Build and validate image
                    // Once done, push to private registry 
                    docker.withRegistry('http://docker-registry:5000') {
                        def hetmanImage = docker.build('hetman')

                        // Make the output directory.
                        sh "mkdir -p output"

                        // Save built image to tar archive
                        sh 'docker save --output hetman.tar hetman'

                        // SBOM generation with anchore/syft
                        sh 'syft hetman.tar -o json=./output/sbom.syft.json -o table' 
                        
                        // Vulnerability scanning with anchore/grype
                        sh 'grype sbom:./output/sbom.syft.json -o json=./output/vuln.grype.json -o table'

                        // Archive validation results as artifacts
                        archiveArtifacts artifacts: 'output/*.json'

                        // Push to local Registry
                        hetmanImage.push()
                    }
                }
            """)
        }
    }
}
