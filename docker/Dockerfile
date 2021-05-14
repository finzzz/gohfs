FROM ubuntu:bionic

# set version label
ARG BUILD_DATE
ARG VERSION
LABEL build_version="Version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="github.com/olofvndrhr"

# gohfs version
ARG GOHFS_VERSION="0.1.3"
ARG GOHFS_ARCH="gohfs-linux-amd64"
ARG GOHFS_PATH="/app/gohfs"

# s6 overlay version
ARG OVERLAY_VERSION="v2.2.0.3"
ARG OVERLAY_ARCH="amd64"

# environment settings
ARG DEBIAN_FRONTEND="noninteractive"
ENV HOME="/root" \
LANGUAGE="en_US.UTF-8" \
LANG="en_US.UTF-8" \
TERM="xterm"

# add s6 overlay
ADD https://github.com/just-containers/s6-overlay/releases/download/${OVERLAY_VERSION}/s6-overlay-${OVERLAY_ARCH}-installer /tmp/
RUN chmod +x /tmp/s6-overlay-${OVERLAY_ARCH}-installer && /tmp/s6-overlay-${OVERLAY_ARCH}-installer / && rm /tmp/s6-overlay-${OVERLAY_ARCH}-installer

# copy sources
COPY sources.list /etc/apt/

# add additional packages
RUN \
 echo "**** install build packages ****" && \
 apt-get update && \
 apt-get install -y --no-install-recommends \
	autoconf \
	apt-utils \
	locales \
	patch \
	tzdata \
	file && \
 echo "**** install runtime packages ****" && \
 apt-get update && \
 apt-get install -y --no-install-recommends \
	git \
	curl \
	wget \
	rsync \
	sudo \
	tar \
	bzip2 \
	zip \
	unzip && \
echo "**** generate locale ****" && \
locale-gen en_US.UTF-8 && \
echo "**** create abc user and make default folders ****" && \
useradd -u 911 -U -d /config -s /bin/false abc && \
usermod -G users abc && \
mkdir -p \
	/app \
	/config \
	/defaults && \
echo "**** cleanup ****" && \
apt-get purge --auto-remove -y && \
apt-get clean && \
rm -rf \
	/tmp/* \
	/var/lib/apt/lists/* \
	/var/tmp/*

# add gohfs
ADD https://github.com/finzzz/gohfs/releases/download/${GOHFS_VERSION}/${GOHFS_ARCH} ${GOHFS_PATH}
RUN chmod +x ${GOHFS_PATH}

# copy local files
COPY root/ /

# ports and volumes
EXPOSE 8080
VOLUME /data

ENTRYPOINT ["/init"]
