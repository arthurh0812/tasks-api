FROM golang:1.16.3 as build

# go environment variables
ENV GO111MODULE="on"
ENV CGO_ENABLED="0"
ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GOFLAGS="-mod=mod"
ENV GOPRIVATE="github.com/arthurh0812"
ENV GOPROXY="https://goproxy.cn,https://gocenter.io,https://goproxy.io,direct"
ENV GOMOD="/build/go.mod"
ENV GOPATH="/home/arthur/go"

WORKDIR /build

# git configuration
RUN apt-get update && apt install git
RUN git config --global url."https://arthurh0812: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@github.com".insteadOf "https://github.com"

COPY go.mod .
COPY go.sum .
RUN go mod download && go mod vendor && go mod verify

COPY . .

RUN go build -o main github.com/arthurh0812/task-app/tasks-api


FROM alpine:3.13.5

WORKDIR /app

COPY --from=build /build/main .

# program-internal environment variables
ENV AUTH_API_SERVICE_HOST="localhost"
ENV TASKS_DIRECTORY="tasks"
ENV PORT='8000'

EXPOSE 8000

RUN mkdir ${TASKS_DIRECTORY}

CMD [ "/app/main" ]