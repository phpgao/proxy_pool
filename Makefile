export PATH := $(GOPATH)/bin:$(PATH)
all: build
build: app
app:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./proxy_pool_darwin_amd64
# 	env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./proxy_pool_linux_386
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./proxy_pool_linux_amd64
# 	env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./proxy_pool_linux_arm
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./proxy_pool_linux_arm64
# 	env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./proxy_pool_windows_386.exe
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./proxy_pool_windows_amd64.exe
# 	env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -o ./proxy_pool_linux_mips64
# 	env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -o ./proxy_pool_linux_mips64le
# 	env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o ./proxy_pool_linux_mips
# 	env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -o ./proxy_pool_linux_mipsle