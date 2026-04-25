.PHONY: build-windows clean run start-db 

clean:
	if exist resource.syso del resource.syso
	if exist versioninfo.go del versioninfo.go

build-windows: clean
	goversioninfo -64 -icon app.ico
	go build -o todo-app.exe ./cmd/api
	@echo Build complete! Check todo-app.exe

start-db:
	docker-compose up -d db

run:
	go run cmd/api/main.go