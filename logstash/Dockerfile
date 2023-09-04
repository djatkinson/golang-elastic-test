
FROM docker.elastic.co/logstash/logstash:7.14.0-arm64
COPY ./logstash/logstash.conf /etc/logstash/conf.d/

CMD logstash -f /etc/logstash/conf.d/logstash.conf