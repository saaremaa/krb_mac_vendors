win: 
	go build -ldflags "-s -w" -o mac.exe
linux: 
	go build -ldflags "-s -w" -o mac
test:
	go test -v -coverprofile cp.out
	go tool cover -html=cp.out