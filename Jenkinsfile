pipeline {
    agent any
    stages {
        stage('🎉🎉🎉 SETUP & BUILD... ') {
            when { 
            not {
              branch 'master'
            }
          }
            steps {
                script {
                    // def root = tool name: '1.8.3', type: 'go'
                    withEnv(["PATH=$PATH:/usr/local/go/bin"]) {
                        sh "go version"
                        sh 'echo "⚠️⚠️⚠️ BUILD SANDEXCARE ⚠️⚠️⚠️"'
                        sh 'rm -rf dist'  
                        sh 'ls -a'  
                        sh 'go version'
                        // sh 'go vet'
                        //sh 'go fmt'
                        sh 'go build -o dist/sandexcare main.go'
                        sh 'sudo systemctl stop api.service'
                        sh 'cp dist/sandexcare /opt/api/'
                        sh 'cp dev.env /opt/api/'
                        sh 'cp -r docs /opt/api/'
                        sh 'cp -r server_websocket /opt/api/'
                        sh 'cd /opt/api/'
                        sh 'sudo systemctl start api.service'
                    }
                }
            }
        }
        stage('🎉🎉🎉 DEPLOY OPEN') {
            when {
                branch "build_dev"
            }
            steps {
                timeout(time: 3, unit: 'MINUTES') {
                    
                }
            }
        }
    }
}
