FROM golang:1.19-bullseye

RUN mkdir /app/src/

WORKDIR /app/src/

COPY go.mod go.sum ./app/src/
RUN go mod download -x

COPY . ./app/src/
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.release=`git rev-parse --short=8 HEAD`'" -o /bin/server ./cmd

EXPOSE 8080

CMD ["go", "run", "./cmd/main.go"]