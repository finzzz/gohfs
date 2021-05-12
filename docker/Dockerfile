FROM ghcr.io/linuxserver/baseimage-alpine:3.13
# set version label
ARG BUILD_DATE
ARG VERSION
LABEL build_version="Version:- ${VERSION} Build-date:- ${BUILD_DATE}"
LABEL maintainer="github.com/olofvndrhr"

# environment settings
ENV GOHFS_VERSION="0.1.2" \
	GOHFS_PATH="/app/gohfs" \
    LD_PRELOAD="/usr/lib/preloadable_libiconv.so"

RUN \
 echo "**** install build packages ****" && \
 apk add --no-cache --virtual=build-dependencies --upgrade \
	autoconf \
	file && \
 echo "**** install runtime packages ****" && \
 apk add --no-cache --upgrade \
	git \
	curl \
	wget \
	rsync \
	sudo \
	tar \
	bzip2 \
	unzip && \
 echo "**** cleanup ****" && \
 apk del --purge \
	build-dependencies && \
 rm -rf \
	/tmp/*

# get gohfs
RUN \
	wget -O "${GOHFS_PATH}" \
		https://github.com/finzzz/gohfs/releases/download/${GOHFS_VERSION}}/gohfs-linux-amd64

# copy local files
COPY root/ /

# ports and volumes
EXPOSE 80
VOLUME /config
