FROM alpine:3.3
MAINTAINER Dylan Meissner "https://github.com/dylanmei/go-nsnitro"

ENV NSNITRO_SERVER ""
ENV NSNITRO_USERNAME ""
ENV NSNITRO_PASSWORD ""

ADD bin/nsnitro /bin/nsnitro
ENTRYPOINT ["/bin/nsnitro"]
