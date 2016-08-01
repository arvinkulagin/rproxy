FROM scratch

ADD rproxy /app
ADD config.json /cfg/config.json

EXPOSE 8888

CMD ["/app", "-config", "/cfg/config.json"]