FROM golang:1.19-bullseye

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.release=`git rev-parse --short=8 HEAD`'" -o /bin/server ./cmd

EXPOSE 8080

CMD ["go", "run", "./cmd/main.go"]