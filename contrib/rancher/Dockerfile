FROM alpine:3.6

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir -p /netapp /etc/netappdvp

RUN mkdir -p /etc/docker
RUN mkdir -p /host/etc/docker/plugins
RUN ln -s /host/etc/docker/plugins /etc/docker/plugins

RUN rmdir /var/run
RUN rmdir /run
RUN mkdir -p /host/run /host/var/run /host/dev/mapper
RUN    ln -s /host/var/run /var/run \
    && ln -s /host/run /run

RUN mkdir -p /host/var/log
RUN rmdir /var/log
RUN ln -s /host/var/log /var/log

RUN mkdir -p /host/var/lib/docker-volumes
RUN ln -s /host/var/lib/docker-volumes /var/lib/docker-volumes

# Use APK mirrors for fault tolerance
RUN printf "http://dl-1.alpinelinux.org/alpine/v3.6/main\nhttp://dl-2.alpinelinux.org/alpine/v3.6/main\nhttp://dl-3.alpinelinux.org/alpine/v3.6/main\nhttp://dl-4.alpinelinux.org/alpine/v3.6/main\nhttp://dl-5.alpinelinux.org/alpine/v3.6/main\n\nhttp://dl-1.alpinelinux.org/alpine/v3.6/community\nhttp://dl-2.alpinelinux.org/alpine/v3.6/community\nhttp://dl-3.alpinelinux.org/alpine/v3.6/community\nhttp://dl-4.alpinelinux.org/alpine/v3.6/community\nhttp://dl-5.alpinelinux.org/alpine/v3.6/community" > /etc/apk/repositories

RUN apk update || true && apk add curl vim screen tmux bash jq

#ADD netappdvp-1.3.tar.gz /netapp
ADD container-launch.sh /netapp
ADD netappdvp /netapp
WORKDIR /etc/netappdvp

#ENV PATH=/netapp:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV PS1="\n[ \w ]\n\# "

ADD chroot-host-wrapper.sh /netapp
RUN    ln -s /netapp/chroot-host-wrapper.sh /netapp/blkid \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/cat \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/df \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/iscsiadm \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/ls \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/lsblk \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/lsscsi \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/mkdir \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/mkfs.ext4 \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/mkfs.xfs \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/mount \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/multipath \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/rescan-scsi-bus \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/rescan-scsi-bus.sh \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/rmdir \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/sh \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/umount \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/stat \
    && ln -s /netapp/chroot-host-wrapper.sh /netapp/sudo \
    && ln -s /netapp/chroot-host-wrapper.sh /sbin/rescan-scsi-bus \
    && ln -s /netapp/chroot-host-wrapper.sh /sbin/rescan-scsi-bus.sh \
    && ln -s /netapp/chroot-host-wrapper.sh /sbin/multipath

#ENTRYPOINT ["/netapp/netappdvp", "--in-container=true"]
CMD ash

