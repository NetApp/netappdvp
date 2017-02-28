#!/usr/bin/env bash

# process debug 
MY_DEBUG=${debug:-false}
case $(echo ${MY_DEBUG} | tr '[:upper:]' '[:lower:]') in
    true)
        DEBUG_SWITCH="--debug=true"
        ;; 
    *)  
        DEBUG_SWITCH="--debug=false"
        ;;  
esac
export DEBUG_SWITCH

# process config
MY_CONFIG=${config:-/etc/netappdvp/config.json}
basefile=$(basename ${MY_CONFIG})
MY_CONFIG="/etc/netappdvp/${basefile}"
export CONFIG_SWITCH="--config=${MY_CONFIG}"

echo Running: /netapp/netappdvp ${DEBUG_SWITCH} ${CONFIG_SWITCH} "${@:1}"
/netapp/netappdvp ${DEBUG_SWITCH} ${CONFIG_SWITCH} "${@:1}"

