FROM golang:1.20 AS build

ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /build/ddns

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /build/. /
ENV GOOGLE_APPLICATION_CREDENTIALS=/sa.json

ENTRYPOINT [ "/ddns" ]

