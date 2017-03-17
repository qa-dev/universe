FROM alpine

RUN mkdir /app
COPY ./universe /app/
COPY ./config.json /app/

EXPOSE 9713

CMD /app/universe
