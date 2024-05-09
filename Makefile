run:
	air

compile:
	echo "Compiling for windows"
	go build -o bin/main-windows-386 main.go

compose:
	docker compose up -d
