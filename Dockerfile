FROM scratch
COPY config /config
COPY kubeaudit /
ENTRYPOINT ["/kubeaudit"]
