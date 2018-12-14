FROM golang:latest

ADD cmd/fast-storage /bin/fast-storage

CMD ["/bin/fast-storage"]