set shell := ["bash", "-cu"]

# 编译 Go CLI 到 dist/proteus
build:
	mkdir -p dist && go build -trimpath -ldflags "-s -w" -o dist/proteus ./cmd/proteus

# 开发模式运行
run:
	go run ./cmd/proteus

# 列出 providers
list:
	go run ./cmd/proteus --list

# 校验配置
validate:
	go run ./cmd/proteus --validate
