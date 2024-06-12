FROM golang:1.22.0 as GO_BUILDER
ARG VER
WORKDIR /server
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o BookQuest ./cmd/main.go

FROM alpine
LABEL org.opencontainers.image.source="https://github.com/robrotheram/bookquest"
WORKDIR /app
ADD app.sample.env /app/app.env
COPY --from=GO_BUILDER /server/BookQuest /app/BookQuest
ADD static /app/static
ADD views /app/views
EXPOSE 8090
ENTRYPOINT ["./BookQuest", "server"]