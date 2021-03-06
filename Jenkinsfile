// Setup milestone to stop previous build from running when a new one is launched
// The result would be:
//  Build 1 runs and creates milestone 1
//  While build 1 is running, suppose build 2 fires. It has milestone 1 and milestone 2. It passes 1, which causes build #1 to abort

def buildNumber = env.BUILD_NUMBER as int
if (buildNumber > 1) milestone(buildNumber - 1)
milestone(buildNumber)

pipeline {
    agent { label 'operator-slave'}
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10')
        timeout(time: 90, unit: 'MINUTES')
    }
    environment {
        WORKING_DIR = "/home/jenkins/go/src/github.com/kiegroup/kogito-cloud-operator/"
        GOPATH = "/home/jenkins/go"
        GOCACHE = "/home/jenkins/go/.cache/go-build"
    }
    stages {
        stage('Clean Workspace') {
            steps {
                dir ("${WORKING_DIR}") {
                    deleteDir()
                }
            }
        }
        stage('Initialize') {
            steps {
                sh "mkdir -p ${WORKING_DIR} && cd ${WORKSPACE} && cp -Rap * ${WORKING_DIR}"
                sh "set +x && oc login --token=\$(oc whoami -t) --server=${OPENSHIFT_API} --insecure-skip-tls-verify"
            }
        }
        stage('Build Kogito Operator') {
            steps {
                dir ("${WORKING_DIR}") {
                    sh """
                        export GOROOT=`go env GOROOT`
                        GO111MODULE=on 
                        go get -u golang.org/x/lint/golint
                        touch /etc/sub{u,g}id
                        usermod --add-subuids 10000-75535 \$(whoami)
                        usermod --add-subgids 10000-75535 \$(whoami)
                        cat /etc/subuid
                        cat /etc/subgid
                        make image_builder=buildah
                    """
                }
            }
            
        }
        stage('Build Kogito CLI') {
            steps {
                dir ("${WORKING_DIR}") {
                    sh "make build-cli"
                }
            }
        }
        stage('Push Operator Image to Openshift Registry') {
            steps {
                dir ("${WORKING_DIR}") {
                    sh """
                        set +x && buildah login -u jenkins -p \$(oc whoami -t) --tls-verify=false ${OPENSHIFT_REGISTRY}
                        cd version/ && TAG_OPERATOR=\$(grep -m 1 'Version =' version.go) && TAG_OPERATOR=\$(echo \${TAG_OPERATOR#*=} | tr -d '"')
                        buildah tag quay.io/kiegroup/kogito-cloud-operator:\${TAG_OPERATOR} ${OPENSHIFT_REGISTRY}/openshift/kogito-cloud-operator:pr-\$(echo \${GIT_COMMIT} | cut -c1-7)
                        buildah push --tls-verify=false docker://${OPENSHIFT_REGISTRY}/openshift/kogito-cloud-operator:pr-\$(echo \${GIT_COMMIT} | cut -c1-7)
                    """
                }
            }
        }
        stage('Running Smoke Testing') {
            steps {
                dir ("${WORKING_DIR}") {
                    sh """
                        make run-smoke-tests load_factor=3 load_default_config=true operator_image=${OPENSHIFT_REGISTRY}/openshift/kogito-cloud-operator operator_tag=pr-\$(echo \${GIT_COMMIT} | cut -c1-7) maven_mirror=${MAVEN_MIRROR_REPOSITORY} concurrent=3
                    """
                }
            }
            post {
                always {
                    dir("${WORKING_DIR}") {
                        archiveArtifacts artifacts: 'test/logs/**/*.log', allowEmptyArchive: true
                        junit testResults: 'test/logs/**/junit.xml', allowEmptyResults: true
                        sh "cd test && go run scripts/prune_namespaces.go"
                    }
                }
            }
        }
    }
}