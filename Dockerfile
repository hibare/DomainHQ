FROM golang:1.21.2-bookworm AS base

# Build main app
FROM base AS build

# Install healthcheck cmd
RUN curl -sfL https://raw.githubusercontent.com/hibare/go-docker-healthcheck/main/install.sh | sh -s -- -d -b /usr/local/bin

WORKDIR /src/

COPY . /src/

# Build DomainHQ
RUN CGO_ENABLED=0 go build -o /bin/domainhq main.go

# Generate final image
FROM scratch

COPY --from=build /bin/domainhq /bin/domainhq

COPY --from=build /usr/local/bin/healthcheck /bin/healthcheck

HEALTHCHECK \
    --interval=30s \
    --timeout=3s \
    CMD ["healthcheck", "--url", "http://localhost:5000/ping/"]

EXPOSE 5000

ENTRYPOINT ["/bin/domainhq"]