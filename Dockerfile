FROM 192.168.0.226/feishuwg/golang:18 AS buildx
#时间时区
RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
WORKDIR /go
COPY . bootstrap-node
WORKDIR /go/bootstrap-node
RUN go env -w GOPROXY=http://192.168.0.227:8081/repository/goproxy/ && \
    go mod download && \
	 go mod tidy -compat=1.18 && \
	 go env -w CGO_ENABLED=0 && go env -w GO111MODULE=on
RUN go build -ldflags "-s -w" -a -o /src/bin/bootstrap .



FROM alpine:3.18
WORKDIR ~
RUN sed -i 's@https://dl-cdn.alpinelinux.org/alpine@http://192.168.0.227:3142/alpine@g' /etc/apk/repositories && \
    apk add --no-cache  curl jq --repository http://mirrors.ustc.edu.cn/alpine/edge/testing

#时间时区
RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
COPY --from=buildx /src/bin /usr/local/bin

CMD ["bootstrap"]