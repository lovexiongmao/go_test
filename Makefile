.PHONY: run build test clean deps

# 运行服务
run:
	go run cmd/server/main.go

# 构建二进制文件
build:
	go build -o bin/server cmd/server/main.go

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/

# 下载依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
vet:
	go vet ./...

# 运行所有检查
check: fmt vet test

