GC=go build
BUILD_NODE_PAR = -trimpath -ldflags "-s -w "
linux:
	GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o bin/site main.go
