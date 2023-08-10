build:
	go build -o bin/ChatterBox.exe main.go

run:
	go run main.go

test:
	go test .

clean:
	go clean .