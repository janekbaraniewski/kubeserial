FROM raspbian/jessie:latest

ENV OPERATOR=/usr/local/bin/kubeserial \
    USER_UID=1001 \
    USER_NAME=kubeserial

# install operator binary
COPY _output/bin/kubeserial ${OPERATOR}

COPY bin/* /usr/local/bin/
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
