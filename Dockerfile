ARG APP_PATH="/go/src/github.com/zerospam/check-firewall"
ARG APP_NAME="firewallChecker"

FROM golang:1.11.1-alpine as builder

ARG APP_PATH
ARG APP_NAME

RUN wget -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN sed -i -e 's/dl-cdn/dl-4/' /etc/apk/repositories && \
    apk --update --no-cache add git

COPY . $APP_PATH
WORKDIR $APP_PATH
RUN dep ensure \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME

FROM scratch

ARG APP_PATH
ARG APP_NAME

COPY --from=builder ${APP_PATH}/${APP_NAME} /${APP_NAME}

CMD ["/firewallChecker"]
