FROM golang:1.11.1-alpine

ARG APP_PATH="/go/src/CheckFirewall"
ARG APP_NAME="firewallChecker"

COPY . $APP_PATH
WORKDIR $APP_PATH
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME

FROM scratch

COPY --from=0 ${APP_PATH}/$APP_NAME /$APP_NAME

CMD ["/$APP_NAME"]
