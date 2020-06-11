FROM golang:1.13-alpine
WORKDIR /go/src/github.com/Shopify/kubeaudit
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kubeaudit -v

FROM scratch
COPY --from=0 /go/src/github.com/Shopify/kubeaudit/kubeaudit .

ENTRYPOINT ["./kubeaudit"]
CMD ["all"]
