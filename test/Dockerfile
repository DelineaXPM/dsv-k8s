FROM dsv-injector:latest

ARG cert_file
ARG key_file
ARG roles_file

COPY ${cert_file} ./dsv.pem
COPY ${key_file} ./dsv.key
COPY ${roles_file} ./roles.json

ENTRYPOINT ["dsv-injector-svc", "-cert", "dsv.pem", "-key", "dsv.key", "-roles", "roles.json" ]
