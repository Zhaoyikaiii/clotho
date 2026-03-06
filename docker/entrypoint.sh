#!/bin/sh
set -e

# First-run: neither config nor workspace exists.
# If config.json is already mounted but workspace is missing we skip onboard to
# avoid the interactive "Overwrite? (y/n)" prompt hanging in a non-TTY container.
if [ ! -d "${HOME}/.clotho/workspace" ] && [ ! -f "${HOME}/.clotho/config.json" ]; then
    clotho onboard
    echo ""
    echo "First-run setup complete."
    echo "Edit ${HOME}/.clotho/config.json (add your API key, etc.) then restart the container."
    exit 0
fi

exec clotho gateway "$@"
