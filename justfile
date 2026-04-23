set shell := ["bash", "-cu"]

# 编译 Go CLI 到 dist/proteus
build-go:
	cd go && mkdir -p ../dist && go build -trimpath -ldflags "-s -w" -o ../dist/proteus ./cmd/proteus
