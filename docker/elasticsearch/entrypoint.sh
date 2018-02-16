#!/usr/local/bin/dumb-init /bin/bash
set -ex

umask 0002

run_as_other_user_if_needed() {
    if [[ "$(id -u)" == "0" ]]; then
        # If running as root, drop to specified UID and run command
        exec chroot --userspec=1000 / "${@}"
    else
        # Either we are running in Openshift with random uid and are a member of the root group
        # or with a custom --user
        exec "${@}"
    fi
}

exec /scripts/makelogs.sh &
run_as_other_user_if_needed "$@"
