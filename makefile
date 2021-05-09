init:
	env GOOS=linux GOARCH=amd64 go build -o releases/gohfs-linux-amd64 cmd/gohfs/main.go
	env GOOS=darwin GOARCH=amd64 go build -o releases/gohfs-macos-amd64 cmd/gohfs/main.go
	env GOOS=windows GOARCH=amd64 go build -o releases/gohfs-amd64.exe cmd/gohfs/main.go

test:
	go run cmd/gohfs/main.go