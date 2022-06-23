FROM golang:1.16
WORKDIR /usr/src/app
ENV GOPROXY=https://goproxy.cn

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./example/main.go ./
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o server
RUN apt-get update && apt-get install -y curl 

EXPOSE 9001
CMD ["/usr/src/app/server"]
