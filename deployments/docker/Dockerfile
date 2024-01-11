# Build the Go binary
FROM golang:1.21 as build_sales-api
# DISABLE CGO
ENV CGO_ENABLED 0
#VERSION CONTROL REFERENCE
ARG BUILD_REF

# Copy the source code into the container
COPY ./main.go ./go.mod /service/

# Build the service binary
WORKDIR /service
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go binary in Alpine
FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=build_sales-api /service /service/
WORKDIR /service
CMD ["./ultimateService"]

#LABEL org.opencontainers.image.created="${BUILD_DATE}"
#      org.opencontainers.image.title="service"
#      org.opencontainers.image.authors="ktruedat"
#      org.opencontainers.image.source="https://github.com/ktruedat/ultimateService"
#      org.opencontainers.image.revision="${BUILD_REF}"
#      org.opencontainers.image.vendor="service"
