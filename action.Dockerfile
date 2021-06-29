# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

FROM buildpack-deps:scm as clone
SHELL ["/bin/bash", "-exo", "pipefail", "-c"]

RUN mkdir actions ;\
	cd actions ;\
	git clone --bare https://github.com/actions/checkout.git ;\
	git -C checkout.git archive --prefix=checkout/ v2 |tar -x ;\
	rm -rf *.git


FROM golang:buster
SHELL ["/bin/bash", "-exo", "pipefail", "-c"]
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update ;\
	apt-get install --no-install-{recommends,suggests} -y \
		apt-transport-https gnupg2 dirmngr ca-certificates ;\
	apt-get clean ;\
	rm -vrf /var/lib/apt/lists/* ;\
	apt-key adv --fetch-keys https://download.docker.com/linux/debian/gpg ;\
	apt-get purge -y gnupg2 dirmngr ;\
	apt-get autoremove --purge -y

ADD action-base.list /etc/apt/sources.list.d/docker.list

RUN apt-get update ;\
	apt-get install --no-install-{recommends,suggests} -y \
		bzip2 docker-ce-cli nodejs upx ;\
	apt-get clean ;\
	rm -vrf /var/lib/apt/lists/*

COPY --from=clone /actions /actions

COPY action.bash Dockerfile /
COPY entrypoint /entrypoint
COPY rootfs /rootfs

CMD ["/action.bash"]
