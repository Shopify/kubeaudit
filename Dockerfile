FROM golang:1.17 AS builder

# no need to include cgo bindings
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# add ca certificates and timezone data files
# hadolint ignore=DL3008
RUN apt-get install --yes --no-install-recommends ca-certificates tzdata

# add unprivileged user
RUN adduser --shell /bin/true --uid 1000 --disabled-login --no-create-home --gecos '' app \
  && sed -i -r "/^(app|root)/!d" /etc/group /etc/passwd \
  && sed -i -r 's#^(.*):[^:]*$#\1:/sbin/nologin#' /etc/passwd

# this is where we build our app
WORKDIR /go/src/app/

# download and cache our dependencies
VOLUME /go/pkg/mod
COPY go.mod go.sum ./
RUN go mod download

# compile kubeaudit
COPY . ./
RUN go build -a -ldflags '-w -s -extldflags "-static"' -o /go/bin/kubeaudit ./cmd/ \
  && chmod +x /go/bin/kubeaudit

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
USER 1000

# entrypoint
ENTRYPOINT ["/kubeaudit"]
CMD ["all"]
