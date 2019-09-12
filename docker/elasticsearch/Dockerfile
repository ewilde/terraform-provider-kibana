ARG ELK_VERSION=6.0.0
ARG ELK_PACK=-oss

FROM docker.elastic.co/elasticsearch/elasticsearch$ELK_PACK:$ELK_VERSION

ARG MAKELOGS_VERSION="makelogs@4.0.3"

USER root
RUN yum install -y openssl wget
RUN yum install -y epel-release && yum install -y nodejs && \
    npm install -g $MAKELOGS_VERSION

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.1.3/dumb-init_1.1.3_amd64 && \
    chmod +x /usr/local/bin/dumb-init

ADD entrypoint.sh /entrypoint.sh
ADD scripts /scripts

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/usr/share/elasticsearch/bin/elasticsearch"]
