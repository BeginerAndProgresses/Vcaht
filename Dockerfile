FROM ubuntu:20.04
COPY v_chat /app/vchat
WORKDIR /app
CMD ["/app/vchat"]