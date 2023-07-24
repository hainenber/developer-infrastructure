FROM jenkins/inbound-agent:jdk17 as jnlp
FROM golang:1.20-alpine

RUN apk -U add git openjdk17-jre docker openrc && rc-update add docker boot

COPY --from=jnlp /usr/local/bin/jenkins-agent /usr/local/bin/jenkins-agent
COPY --from=jnlp /usr/share/jenkins/agent.jar /usr/share/jenkins/agent.jar

ENTRYPOINT ["/usr/local/bin/jenkins-agent"]