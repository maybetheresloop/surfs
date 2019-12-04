FROM golang:1.13.4-alpine3.10

COPY bin/surfs-block /surfs-block
COPY bin/surfs-meta /surfs-meta

RUN mkdir /data

EXPOSE 5678 5679

ENTRYPOINT ["/surfs-block"]
CMD ["--datadir", "/data", "--port", "5678", "-VV"]
