build-and-push-agent: 
	docker build --file ./jenkins/jenkins.agent.Dockerfile --tag jenkins-agent-for-golang .
	docker tag jenkins-agent-for-golang localhost:5000/jenkins-agent-for-golang 
	docker push localhost:5000/jenkins-agent-for-golang 

build-jcasc-config:
	cd ./scripts && go run ./...

run: build-jcasc-config
	docker-compose --file ./docker-compose.yaml up --detach --remove-orphans --build
	$(MAKE) build-and-push-agent

stop:
	docker-compose --file ./docker-compose.yaml down -v