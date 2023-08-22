FROM golang:alpine 

RUN mkdir /app
WORKDIR /app
COPY go.mod .

RUN go mod tidy
COPY . .
RUN go build cmd/main.go

EXPOSE 8081
CMD [ "./main" ]