FROM alpine

LABEL version="0.0.3"
EXPOSE 1080

RUN mkdir /app
WORKDIR  /app
VOLUME [ "/app" ]

ADD erproxy /bin
RUN chmod +x /bin/erproxy
ADD config.yml /
ADD enter.sh /
RUN chmod +x /enter.sh

CMD ["sh","/enter.sh"]

