FROM fedora:latest as builder


RUN dnf install -y golang git
RUN dnf update -y && dnf upgrade -y

RUN git clone https://github.com/aaraney/ht.git && cd ht && go build .

FROM scratch
COPY --from=builder /ht/ht /ht
CMD ["/ht"]
