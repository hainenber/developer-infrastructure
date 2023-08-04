build-and-push-agent: 
	docker build --file ./jenkins/jenkins.agent.Dockerfile --tag jenkins-agent-for-golang ./jenkins
	docker tag jenkins-agent-for-golang localhost:32001/jenkins-agent-for-golang 
	docker push localhost:32001/jenkins-agent-for-golang 

build-jcasc-config:
	cd ./scripts/build-jenkins-jobs-into-jcasc-config && go run ./...

add-athens-host-as-jenkins-global-var:
	cd ./scripts/add-athens-host-as-jenkins-global-var && go run ./...

start-k8s: 
	$(MAKE) build-jcasc-config
	docker build --file ./jenkins/jenkins.server.Dockerfile --tag jenkins-server ./jenkins
	# sleep 5
	# $(MAKE) build-and-push-agent
	cp ./jenkins/casc-configs/jcasc.yaml ./deploy/kubernetes/jcasc.yaml
	cp ./jenkins/.secrets ./deploy/kubernetes/.secrets
	kubectl apply -f ./deploy/kubernetes/namespace.yml
	kubectl apply -k ./deploy/kubernetes/
	kubectl apply -f ./deploy/kubernetes/jenkins/ 
	kubectl apply -f ./deploy/kubernetes/docker-registry/ 

stop-k8s:
	kubectl delete -k ./deploy/kubernetes/ || true
	kubectl delete -f ./deploy/kubernetes/jenkins/ || true
	kubectl delete -f ./deploy/kubernetes/docker-registry/ || true

start: 
	docker-compose --file ./deploy/docker-compose/docker-compose.yaml up --detach --remove-orphans --build
	$(MAKE) add-athens-host-as-jenkins-global-var
	$(MAKE) build-jcasc-config
	$(MAKE) build-and-push-agent

stop:
	docker-compose --file ./deploy/docker-compose/docker-compose.yaml down -v --remove-orphans

restart:
	$(MAKE) stop && $(MAKE) start

restart-k8s:
	$(MAKE) stop-k8s && $(MAKE) start-k8s