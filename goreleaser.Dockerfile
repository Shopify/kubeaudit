FROM scratch

COPY kubeaudit /

ENTRYPOINT ["/kubeaudit"]
CMD ["all"]
