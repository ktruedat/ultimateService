# Build the Go binary
FROM golang:1.21 as build_sales-api
# DISABLE CGO
ENV CGO_ENABLED 0
#VERSION CONTROL REFERENCE
ARG BUILD_REF

# Copy the source code into the container
COPY . /service/

# Build the service binary
WORKDIR /service/app/services/sales-api
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go binary in Alpine
FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=build_sales-api /service/app/services/sales-api/sales-api /service/sales-api
WORKDIR /service
CMD ["./sales-api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="ktruedat" \
      org.opencontainers.image.source="https://github.com/ktruedat/ultimateService" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="ktruedat"
