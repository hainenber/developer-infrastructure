FROM jenkins/inbound-agent:jdk17 as jnlp
FROM anchore/syft:latest as syft
FROM anchore/grype:latest as grype
FROM golang:1.20-alpine as golang

# Install prerequisites
RUN apk -U add git openjdk17-jre docker openrc && \
    rc-update add docker boot

COPY --from=syft syft /usr/local/bin/syft
COPY --from=grype grype /usr/local/bin/grype
COPY --from=jnlp /usr/local/bin/jenkins-agent /usr/local/bin/jenkins-agent
COPY --from=jnlp /usr/share/jenkins/agent.jar /usr/share/jenkins/agent.jar

ENTRYPOINT ["/usr/local/bin/jenkins-agent"]