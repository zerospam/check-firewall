ARG APP_PACKAGE="github.com/zerospam/check-firewall"
ARG APP_PATH="/go/src/${APP_PACKAGE}"
ARG APP_NAME="firewallChecker"

FROM zerospam/go-dep-docker as builder

ENV CGO_ENABLED=0
ENV GOOS=linux

ARG APP_PACKAGE
ARG APP_PATH
ARG APP_NAME

COPY . $APP_PATH
WORKDIR $APP_PATH

RUN dep ensure \
    && go test -v ${APP_PACKAGE}/test

RUN go build -a -installsuffix cgo -o $APP_NAME

FROM scratch

ARG APP_PATH
ARG APP_NAME

COPY --from=builder ${APP_PATH}/${APP_NAME} /${APP_NAME}

CMD ["/firewallChecker"]
