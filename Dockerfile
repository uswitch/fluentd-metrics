FROM scratch
LABEL maintainer Tom Taylor <tom.taylor@uswitch.com>

ADD fluentd-metrics /
ENTRYPOINT ["/fluentd-metrics"]
