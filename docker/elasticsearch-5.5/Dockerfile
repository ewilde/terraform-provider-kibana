FROM docker.elastic.co/elasticsearch/elasticsearch:5.5.3

USER root
RUN yum install -y epel-release
RUN yum install -y openssl
RUN yum install -y nodejs
RUN npm install -g makelogs

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.1.3/dumb-init_1.1.3_amd64 && \
    chmod +x /usr/local/bin/dumb-init

USER elasticsearch

ADD entrypoint.sh /entrypoint.sh
ADD scripts /scripts

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/usr/share/elasticsearch/bin/elasticsearch"]
