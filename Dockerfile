FROM golang:1.15.0-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# add ca certificates and timezone data files
# hadolint ignore=DL3018
RUN apk add -U --no-cache ca-certificates tzdata

# add unprivileged user
RUN adduser -s /bin/true -u 1000 -D -h /app app \
  && sed -i -r "/^(app|root)/!d" /etc/group /etc/passwd \
  && sed -i -r 's#^(.*):[^:]*$#\1:/sbin/nologin#' /etc/passwd

WORKDIR /go/src/app/

COPY . ./
RUN go build -ldflags '-w -s -extldflags "-static"' -o /go/bin/kubeaudit -v \

#
# ---
#

# start with empty image
FROM scratch

# add-in our timezone data file
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# add-in our unprivileged user
COPY --from=builder /etc/passwd /etc/group /etc/shadow /etc/

# add-in our ca certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# add-in our application
COPY --from=builder --chown=app /go/bin/kubeaudit /kubeaudit

# from now on, run as the unprivileged user
USER app

# entrypoint
ENTRYPOINT ["/kubeaudit"]
