FROM bianjie/irita:latest

COPY . /scripts

RUN sh /scripts/setup.sh

EXPOSE 26657
EXPOSE 9090

CMD irita start