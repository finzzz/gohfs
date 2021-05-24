init:
	env GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o releases/gohfs-linux-amd64 cmd/gohfs/main.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -o releases/gohfs-macos-amd64 cmd/gohfs/main.go
	env GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o releases/gohfs-amd64.exe cmd/gohfs/main.go

test:
	go run cmd/gohfs/main.go

compress:
	upx --ultra-brute releases/gohfs-linux-amd64
	upx --ultra-brute releases/gohfs-macos-amd64
	cp releases/gohfs-amd64.exe releases/gohfs-amd64-packed.exe
	upx --ultra-brute releases/gohfs-amd64-packed.exe