FROM golang:1.10.3-stretch
COPY vendor /go/src/github.com/ww-tech/aws-sts-proxy/vendor
COPY helpers /go/src/github.com/ww-tech/aws-sts-proxy/helpers
COPY main.go /go/src/github.com/ww-tech/aws-sts-proxy/
COPY Gopkg.lock /go/src/github.com/ww-tech/aws-sts-proxy/
COPY Gopkg.toml /go/src/github.com/ww-tech/aws-sts-proxy/
COPY .aws /root/.aws
WORKDIR /go/src/github.com/ww-tech/aws-sts-proxy/
RUN go build -o app 
CMD ["./app"]