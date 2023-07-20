build-agent: 
	docker build --file ./jenkins/jenkins.agent.Dockerfile --tag jenkins-agent-for-golang .
	docker tag jenkins-agent-for-golang localhost:5000/jenkins-agent-for-golang 
	docker push localhost:5000/jenkins-agent-for-golang 

run: 
	docker-compose --file ./docker-compose.yaml up --detach --remove-orphans

stop:
	docker-compose --file ./docker-compose.yaml down -v