FROM jenkins/jenkins:2.413-jdk17

# Use root to install prerequisites
USER root

RUN apt-get update && apt-get install -y lsb-release

RUN curl -fsSLo /usr/share/keyrings/docker-archive-keyring.asc \
  https://download.docker.com/linux/debian/gpg

RUN echo "deb [arch=$(dpkg --print-architecture) \
  signed-by=/usr/share/keyrings/docker-archive-keyring.asc] \
  https://download.docker.com/linux/debian \
  $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list
RUN apt-get update && apt-get install -y docker-ce-cli

# Store recommended Jenkins plugins into Jenkins home dir
COPY plugins.txt /var/jenkins_home/plugins.txt

# Run Jenkins with non-root "jenkins" user
USER jenkins

# Install plugins 
RUN jenkins-plugin-cli --plugin-file /var/jenkins_home/plugins.txt
