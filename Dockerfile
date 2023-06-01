# build stage
FROM golang:1.20 AS build

ENV GO111MODULE=on

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -x -ldflags "-s -w" -o main .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=build /app/main .

# Define your environment variables here
ENV TZ=Europe/Paris \
    PORT=8080 \
    LOG_FILE=/app/logs/golang-app.log \
    KEYCLOAK_URL=http://localhost:8081/auth/realms/ \
    REALM=local \
    CLIENT_ID=gin \
    CLIENT_SECRET=secret


# Expose the port on which the application will listen
EXPOSE $PORT

VOLUME /var/logs/app

USER 1000:100

# Start the application
ENTRYPOINT ["/app/main"]
