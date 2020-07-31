# Creates .env template if nonexistent and shows steps to completes before starting the container
init:
	./share/scripts/init-message.py
# Starts docker container
run:
	docker-compose down && \
	docker-compose up
# Starts docker container
run-detached:
	docker-compose down && \
	docker-compose up -d
# Starts docker container
run-build:
	docker-compose down && \
	docker-compose up --build
# Starts docker container
run-build-d:
	docker-compose down && \
	docker-compose up --build -d