FROM bianjieai/cschain:latest

COPY . /scripts

RUN sh /scripts/setup.sh

EXPOSE 26657

CMD cschain start