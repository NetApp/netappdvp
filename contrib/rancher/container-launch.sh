#!/usr/bin/env bash

VOLUME_DRIVER_NAME="netapp"
VOLUME_DRIVER_NAME_SET_FROM_CLI=false
CONFIG_LOCATION="/etc/netappdvp/config.json"
CONFIG_JSON="${CONFIG_LOCATION}"
for i in "$@"; do
	case $i in
		--volume-driver=*)
                        VOLUME_DRIVER_NAME_SET_FROM_CLI=true
			VOLUME_DRIVER_NAME="${i#*=}"
			shift 
		;;
		--config=*)
			CONFIG_LOCATION="${i#*=}"
			shift 
		;;
		*)
		    # unknown option, ignore
		;;
	esac
done
printf "VOLUME_DRIVER_NAME = ${VOLUME_DRIVER_NAME}\n"
printf "CONFIG_LOCATION = ${CONFIG_LOCATION}\n"

case $CONFIG_LOCATION in
    http://*|https://*)
        printf "Looks like a URL!\n"
        URL="${CONFIG_LOCATION}"
        ;;
    *)
        if [ ! -e "${CONFIG_LOCATION}" ]; then
            printf "Could not find config json ${CONFIG_LOCATION}\n"
            exit 1
        fi
        CONFIG_JSON="${CONFIG_LOCATION}" # they gave a file, not a URL, use that
        ;;
esac

# ensure the directory where we store volume mounts exists
HOST_VOL_DIR="/host/var/lib/docker-volumes"
if [ ! -d "${HOST_VOL_DIR}" ]; then

    # don't see it, let's create the directory
    printf "Creating missing ${HOST_VOL_DIR}\n"
    mkdir -p ${HOST_VOL_DIR}

    # problem creating the directory
    if [ $? -ne 0 ]; then
        printf "Could not create missing ${HOST_VOL_DIR}\n"
        exit 1
    fi

    # no error code and the directory still doesn't exist, something went wrong let's error out
    if [ ! -d "${HOST_VOL_DIR}" ]; then
        printf "Could not create missing ${HOST_VOL_DIR}\n"
        exit 1
    fi
else
    printf "Using existing ${HOST_VOL_DIR}\n"
fi


# if CONFIG_JSON exists, use it;  otherwise, download and filter through jq
if [ ! -e "${CONFIG_JSON}" ]; then
    RANCHER_JSON="/etc/netappdvp/rancher-config.json"

    rm -f ${CONFIG_JSON}
    curl --connect-timeout 5 --header 'Accept: application/json' "${URL}" -o ${RANCHER_JSON}
    if [ $? -ne 0 ]; then
        printf "Could not download rancher json from ${URL}\n"
        exit 1
    fi

    if [ ! -e "${RANCHER_JSON}" ]; then
        printf "Could not find downloaded rancher json at ${RANCHER_JSON}\n"
        exit 1
    fi

    if [ "$VOLUME_DRIVER_NAME_SET_FROM_CLI" = false ] ; then
        VOLUME_DRIVER_NAME=$(cat ${RANCHER_JSON} | jq -r '.volumeDriverName')
        printf "Using volumeDriverName=${VOLUME_DRIVER_NAME} (from Rancher JSON)\n"
    else
        printf "Using volumeDriverName=${VOLUME_DRIVER_NAME} (from CLI)\n"
    fi

    # filter it, expanding the "Types" entry if it is a string instead of an array
    cat ${RANCHER_JSON} \
      | jq 'with_entries( if (.key == "Types" and (.value | type) == "string" and .value != "" ) then (.value |= fromjson) else (.) end )' \
      | jq 'with_entries( if (.key == "DefaultVolSz" and (.value | type) == "string" and .value != "" ) then (.value |= tonumber) else (.) end)' \
      | jq 'with_entries( if (.key == "webProxyVerifyTLS" and (.value | type) == "string" and .value != "" ) then (.value |= startswith("true")) else (.) end)' \
      | jq 'with_entries( if (.key == "webProxyUseHTTP" and (.value | type) == "string" and .value != "" ) then (.value |= startswith("true")) else (.) end)' \
      | jq 'with_entries( if (.key == "debug" and (.value | type) == "string" and .value != "" ) then (.value |= startswith("true")) else (.) end)' \
      > ${CONFIG_JSON}
    #  | tee ${CONFIG_JSON}

    # non zero exit code, there was a problem filtering the rancher json response
    if [ $? -ne 0 ]; then
        printf "Problem filtering rancher json\n"
        exit 1
    fi

    # could not find the filtered ${CONFIG_JSON} file
    if [ ! -e "${CONFIG_JSON}" ]; then
        printf "Could not find config json ${CONFIG_JSON}\n"
        exit 1
    else
        printf "Using config ${CONFIG_JSON}\n"
    fi

else
    printf "Using existing config ${CONFIG_JSON}\n"
fi

OUTPUT=$(/usr/bin/env -i PATH='/host/sbin:/host/bin:/host/usr/bin' /usr/bin/which rescan-scsi-bus)
if [ $? -ne 0 ]; then
    rm -f /netapp/rescan-scsi-bus /sbin/rescan-scsi-bus
fi

OUTPUT2=$(/usr/bin/env -i PATH='/host/sbin:/host/bin:/host/usr/bin' /usr/bin/which rescan-scsi-bus.sh)
if [ $? -ne 0 ]; then
    rm -f /netapp/rescan-scsi-bus.sh /sbin/rescan-scsi-bus.sh
fi

if [ ! -e "/host/sbin/multipath" ]; then
    rm -f /netapp/multipath /sbin/multipath
fi

if [ ! -e "/host/var/lib/docker-volumes" ]; then
    mkdir -p /host/var/lib/docker-volumes
fi

/usr/bin/env -i PATH='/netapp:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin' /netapp/netappdvp --volume-driver=${VOLUME_DRIVER_NAME} --config=${CONFIG_JSON}
