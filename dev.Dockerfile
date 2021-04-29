FROM golang:1.16.3-alpine

# go environment variables
ENV GO111MODULE="on"
ENV CGO_ENABLED="0"
ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GOFLAGS="-mod=mod"
ENV GOPRIVATE="github.com/arthurh0812"
ENV GOPROXY="https://goproxy.cn,https://gocenter.io,https://goproxy.io,direct"
ENV GOMOD="/app/go.mod"
ENV GOPATH="/home/arthur/go"

RUN apk add git && apk update && apk upgrade
RUN git config --global url."https://arthurh0812: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@github.com".insteadOf "https://github.com"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod vendor && go mod verify

COPY . .

ENV AUTH_API_SERVICE_HOST="localhost:80"
ENV TASKS_FOLDER="tasks"
ENV PORT='8000'

EXPOSE 8000

CMD [ "go", "run", "github.com/arthurh0812/task-app/tasks-api" ]