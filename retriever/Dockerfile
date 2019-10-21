FROM golang:1.12

WORKDIR /go/src/app

# we need to exit 0 cause mongo-driver install return a non zero value
RUN go get go.mongodb.org/mongo-driver; exit 0
RUN go get gopkg.in/mgo.v2/bson
COPY . .
RUN go build -o main
CMD ["./main"]
