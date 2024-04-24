FROM registry.cn-hangzhou.aliyuncs.com/godev/golang:1.22 AS builder

LABEL stage=gobuilder

WORKDIR /src
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64 \
GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod download
# 执行构建
COPY . .
RUN make build

FROM registry.cn-hangzhou.aliyuncs.com/godev/alpine:latest

WORKDIR /app
EXPOSE 80

COPY --from=builder /src/dist/mflow /app/mflow-api
COPY --from=builder /src/etc /app/etc

# 默认配置
ENV APP_NAME=mflow \
APP_DOMAIN=console.mdev.group \
HTTP_HOST=127.0.0.1 \
HTTP_PORT=8090 \
GRPC_HOST=127.0.0.1 \
GRPC_PORT=18090 \
MONGO_ENDPOINTS=127.0.0.1:27017 \
MONGO_DATABASE=mflow \
MCENTER_GRPC_ADDRESS=127.0.0.1:18010 \
MCENTER_CLINET_ID=mflow \
MCENTER_CLIENT_SECRET=mflow

CMD ["./mflow-api", "start", "-t", "env"]