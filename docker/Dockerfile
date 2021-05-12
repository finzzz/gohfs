FROM ghcr.io/linuxserver/baseimage-ubuntu:bionic
# set version label
ARG BUILD_DATE
ARG VERSION
LABEL build_version="Version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="github.com/olofvndrhr"

# gohfs version
ENV GOHFS_VERSION="0.1.2" \
	GOHFS_PATH="/app/gohfs" \
	GOHFS_ARCH="gohfs-linux-amd64"

# environment settings
ARG DEBIAN_FRONTEND="noninteractive"
ENV HOME="/root" \
	LANGUAGE="en_US.UTF-8" \
	LANG="en_US.UTF-8" \
	TERM="xterm"

RUN \
 echo "**** install build packages ****" && \
 apt-get update && \
 apt-get install -y --no-install-recommends \
	autoconf \
	file && \
 echo "**** install runtime packages ****" && \
 apt-get update && \
 apt-get install -y --no-install-recommends \
	cron \
	git \
	curl \
	wget \
	rsync \
	sudo \
	tar \
	bzip2 \
	unzip && \
 echo "**** cleanup ****" && \
 apt-get purge --auto-remove -y && \
 apt-get clean && \
 rm -rf \
	/tmp/* \
	/var/lib/apt/lists/* \
	/var/tmp/*

# get gohfs
RUN \
	wget -O "${GOHFS_PATH}" \
		https://github.com/finzzz/gohfs/releases/download/"${GOHFS_VERSION}"/"${GOHFS_ARCH}"

# copy local files
COPY root/ /

# ports and volumes
EXPOSE 80
VOLUME /config
