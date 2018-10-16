def label = "replicator-${UUID.randomUUID().toString()}"
podTemplate(label: label, inheritFrom: 'default', containers: [
    containerTemplate(name: 'docker', image: 'docker', ttyEnabled: true, command: 'cat')
 ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
  ]) {
    def image="zerospam/check-firewall"
    def tag = "0.1"
    def builtImage = null

    node (label) {
        gitInfo = checkout scm
         container('docker') {
            stage('docker build') {
               builtImage = docker.build("${image}:${env.BUILD_ID}")
            }

            if (gitInfo.GIT_BRANCH.equals('master')) {
                // master branch release
                stage('Push docker image to Docker Hub') {
                    docker.withRegistry('https://index.docker.io/v1/', 'docker-hub') {
                            builtImage.push('latest')
                            builtImage.push("${tag}")
                    }
                } // stage
            } // if master branch

        } //container
    } //node
} //pipeline