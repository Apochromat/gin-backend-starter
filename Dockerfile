FROM golang:1.23-bullseye AS build-stage

WORKDIR /app

COPY ./ ./

RUN go mod tidy

RUN go build -o /executable

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /executable /executable

EXPOSE 4000

USER nonroot:nonroot

ENTRYPOINT ["/executable"]