FROM dsv-injector:latest

ARG cert_file
ARG key_file
ARG roles_file

WORKDIR /home/dsv

COPY --chown=dsv:dsv ${cert_file} ./dsv.pem
COPY --chown=dsv:dsv ${key_file} ./dsv.key
COPY --chown=dsv:dsv ${roles_file} ./roles.json

ENTRYPOINT ["dsv-injector-svc", "-cert", "dsv.pem", "-key", "dsv.key", "-roles", "roles.json" ]
