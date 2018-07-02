FROM arnaudgiuliani/golang-glide as builder

ENV APP_PATH=/go/src/github.com/Turbots/go-encrypt

RUN mkdir -p $APP_PATH
#ADD . $APP_PATH
WORKDIR $APP_PATH

RUN glide install -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder $APP_PATH /app/
WORKDIR /app
CMD ["./main"]