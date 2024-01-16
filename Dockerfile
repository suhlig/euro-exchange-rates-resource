FROM golang as build
WORKDIR /usr/local/src/resource
COPY . .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/resource -ldflags '-extldflags "-static"'

FROM registry.access.redhat.com/ubi8-minimal:latest
RUN mkdir -p /opt/resource
COPY --from=build /usr/local/bin/resource /usr/local/bin/
RUN    printf '#!/usr/bin/env bash\n/usr/local/bin/resource check "$@"' > /opt/resource/check \
    && printf '#!/usr/bin/env bash\n/usr/local/bin/resource put "$@"' > /opt/resource/out \
    && printf '#!/usr/bin/env bash\n/usr/local/bin/resource get "$@"' > /opt/resource/in \
    && chmod +x /opt/resource/*
