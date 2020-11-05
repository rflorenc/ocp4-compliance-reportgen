FROM registry.access.redhat.com/ubi8

USER 0

RUN dnf copr enable -y openscapmaint/openscap-latest && \
    dnf install -y nginx bzip2 procps scap-security-guide && \
    dnf clean all && \
    mkdir -p /opt/nginx/html

ENV HOME=/opt/nginx

EXPOSE 8080
VOLUME /opt/nginx/html

ADD nginx.conf /opt/nginx/nginx.conf
ADD reportgen /opt/nginx/
ADD css /opt/nginx/css

RUN chmod -R 775 /opt/nginx

CMD /opt/nginx/reportgen

USER 1001