FROM alpine:3.17

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apl/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose.exe
RUN chmod +x /bin/goose.exe

WORKDIR /root

ADD migrations/*.sql migrations/
ADD migrations.sh .
ADD .env .

RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"]