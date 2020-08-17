FROM alpine:3.12
COPY scheriff /usr/local/bin/scheriff
ENTRYPOINT ["/usr/local/bin/scheriff"]
