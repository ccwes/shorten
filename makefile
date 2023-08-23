# Binary name
BINARY=drone-webhook
GOBUILD=go build -o ${BINARY}
GOCLEAN=go clean
VERSION=0.0.2

# Build
build:
	$(GOCLEAN)
	$(GOBUILD)


clean:
	$(GOCLEAN)
	$(RMTARGZ)


release:
	# Clean
	$(GOCLEAN)
	# Build for mac
	#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD)
	# Build for arm
#	$(GOCLEAN)
	#CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD)
	# Build for linux
#	$(GOCLEAN)

	GOOS=linux GOARCH=amd64 $(GOBUILD)
	# Build for win
#	$(GOCLEAN)
#	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD).exe
#	$(GOCLEAN)

#windows64
GOOS=windows GOARCH=amd64 go build -o bin/app-amd64.exe app.go
#window32
GOOS=windows GOARCH=386 go build -o bin/app-386.exe app.go
# macos 64-bit
GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin app.go
# macos 32-bit
GOOS=darwin GOARCH=386 go build -o bin/app-386-darwin app.go

# linux 64-bit
GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux app.go

# linux 32-bit
GOOS=linux GOARCH=386 go build -o bin/app-386-linux app.go
docker:
	GOOS=linux GOARCH=amd64 $(GOBUILD)
	docker build -t drone/webhook . --tag registry.cn-shanghai.aliyuncs.com/zzf2001/drone-webhook
	docker push registry.cn-shanghai.aliyuncs.com/zzf2001/drone-webhook