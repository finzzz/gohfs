init:
	env GOOS=linux GOARCH=amd64 go build -o releases/gohfs-linux-amd64 gohfs.go
	env GOOS=darwin GOARCH=amd64 go build -o releases/gohfs-macos-amd64 gohfs.go
	env GOOS=windows GOARCH=amd64 go build -o releases/gohfs-amd64.exe gohfs.go
