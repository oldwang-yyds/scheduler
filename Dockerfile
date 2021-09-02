FROM golang:1.16
WORKDIR $GOPATH/src/git.ucloudadmin.com/kun/scheduler
ADD . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o scheduler main.go
CMD ["scheduler"]