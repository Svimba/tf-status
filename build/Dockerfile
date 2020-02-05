FROM docker-dev-local.docker.mirantis.net/tungsten/contrail-status:5.1-dev

ENV OPERATOR=/usr/local/bin/tf-status-proxy
    #  USER_UID=1001 \
    #  USER_NAME=tungstenfabric-operator

# install operator binary
COPY build/_output/bin/tf-status-proxy ${OPERATOR}
COPY build/bin/entrypoint /
RUN chmod +x /entrypoint

# COPY build/bin /usr/local/bin
# RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/entrypoint"]

# USER ${USER_UID}