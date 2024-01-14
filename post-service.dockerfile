FROM alpine:latest

RUN mkdir /app

COPY post-service /app

CMD [ "/app/post-service"]