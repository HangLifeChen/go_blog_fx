FROM golang:1.23 AS builder

# set environment variables
ENV GO111MODULE=on \
    TZ=Asia/Shanghai \
    GOPROXY=https://goproxy.cn,direct

# set working directory
WORKDIR /app

# copy go modules and dependencies
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# copy project files
COPY . .

# build binary file
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server -v -ldflags "-s -w" ./cmd/server

# use Alpine 3.19 as runtime environment
FROM alpine:3.19

# install timezone data and set timezone to Shanghai
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# define BUILD_MODE argument
ARG BUILD_MODE

# copy compiled binary file
COPY --from=builder /app/server /server
COPY --from=builder /app/i18n /i18n
COPY --from=builder /app/conf/config.${BUILD_MODE}.yaml /conf/config.yaml

# give execution permission
RUN chmod +x /server

# set container start command
CMD ["/server", "-c", "/conf/config.yaml"]
