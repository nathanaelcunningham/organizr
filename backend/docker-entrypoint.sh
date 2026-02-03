#!/bin/sh
set -e

# Update user/group IDs if PUID/PGID are set
if [ -n "$PUID" ] && [ -n "$PGID" ]; then
    echo "Setting user app to UID=$PUID and GID=$PGID"

    # Modify group ID if different
    if [ "$(id -g app)" != "$PGID" ]; then
        groupmod -o -g "$PGID" app
    fi

    # Modify user ID if different
    if [ "$(id -u app)" != "$PUID" ]; then
        usermod -o -u "$PUID" app
    fi
fi

# Fix ownership of data directory and app files
chown -R app:app /data /app

# Drop privileges and run the command as the app user
exec su-exec app "$@"
