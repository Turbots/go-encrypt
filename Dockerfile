FROM openshift/golang as builder

ENV GOBIN /go/bin
RUN mkdir /app
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get -u github.com/golang/dep/...
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /app/main .

FROM scratch

COPY --from=builder /app/main /app/
WORKDIR /app
CMD ["./main"]