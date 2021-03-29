# builder image
FROM golang:1.16-alpine as builder
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build 



# final stage
FROM alpine:latest

# Create new user, avoid running processes as root
RUN apk update && \
    addgroup -S gollp && adduser -S gollp -G gollp
USER gollp

# Copy built binaries from builder. Make gollp user the owner
COPY --from=builder --chown=gollp /build/gollp /app/
WORKDIR /app
ENTRYPOINT [ "./gollp" ]