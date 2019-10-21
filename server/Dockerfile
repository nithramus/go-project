FROM golang:1.12

WORKDIR /go/src/app
RUN go get go.mongodb.org/mongo-driver; exit 0
RUN go get gopkg.in/mgo.v2/bson
COPY . .
RUN go build -o main
EXPOSE 80
CMD ["./main"]
