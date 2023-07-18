run: 
	docker-compose --file ./docker-compose.yaml up --detach --remove-orphans

stop:
	docker-compose --file ./docker-compose.yaml down -v