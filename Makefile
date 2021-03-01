
build:
	go build -o bin/castlebot cmd/castlebot/main.go

docker: build
	docker build . -t docker.local.pw10n.pw/castlecoders/castlebot

push: docker 
	docker push docker.local.pw10n.pw/castlecoders/castlebot