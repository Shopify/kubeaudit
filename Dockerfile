FROM golang:1.15.0-alpine AS builder

# no need to include cgo bindings
ENV CGO_ENABLED=0

# add ca certificates and timezone data files
# hadolint ignore=DL3018
RUN apk add -U --no-cache ca-certificates tzdata

# add unprivileged user
RUN adduser -s /bin/true -u 1000 -D -h /app app \
  && sed -i -r "/^(app|root)/!d" /etc/group /etc/passwd \
  && sed -i -r 's#^(.*):[^:]*$#\1:/sbin/nologin#' /etc/passwd

# this is where we build our app
WORKDIR /go/src/app/

# download and cache our dependencies
COPY go.mod go.sum ./
RUN go mod download

# compile kubeaudit
COPY . ./
RUN go build -ldflags '-w -s -extldflags "-static"' -o /go/bin/kubeaudit -v \
  && chown +x /go/bin/kubeaudit

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
