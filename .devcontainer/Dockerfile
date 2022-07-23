# Note: You can use any Debian/Ubuntu based image you want.
FROM mcr.microsoft.com/vscode/devcontainers/universal:2-focal


# Enable new "BUILDKIT" mode for Docker CLI
ENV DOCKER_BUILDKIT=1

# Options
ARG INSTALL_ZSH="true"
ARG UPGRADE_PACKAGES="false"
ARG USE_MOBY="true"
ARG USERNAME=codespace
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG ENABLE_NONROOT_DOCKER="true"
ARG DOCKER_VERSION="latest"
# Install needed packages and setup non-root user. Use a separate RUN statement to add your own dependencies.
USER root
COPY library-scripts/*.sh /tmp/library-scripts/
RUN apt-get update && /bin/bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" "true" "true" \
    && /bin/bash /tmp/library-scripts/docker-in-docker-debian.sh "${ENABLE_NONROOT_DOCKER}" "${USERNAME}" "${USE_MOBY}" "${DOCKER_VERSION}" \
    && /bin/bash /tmp/library-scripts/kubectl-helm-debian.sh "latest" "latest" "latest" \
    && mkdir -p /home/${USERNAME}/.minikube \
    && chown ${USERNAME} /home/${USERNAME}/.minikube \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts/
VOLUME [ "/var/lib/docker" ]
ENTRYPOINT ["/usr/local/share/docker-init.sh"]
CMD ["sleep", "infinity"]
RUN $(curl -fSSl https://github.com/wercker/stern/releases/download/1.11.0/stern_linux_amd64 -o ./stern) && chmod +x ./stern && mv ./stern /usr/local/bin/ && echo 'stern log tool installed'
ENV GOROOT=
ENV GOPATH=/home/${USERNAME}/go
ENV PATH=/home/${USERNAME}/go/bin:$PATH

USER ${USERNAME}


# Setting the ENTRYPOINT to docker-init.sh will start up the Docker Engine
# inside the container "overrideCommand": false is set in devcontainer.json.
# The script will also execute CMD if you need to alter startup behaviors.


# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>
