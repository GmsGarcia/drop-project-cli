#!/bin/bash

# Must be run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root." >&2
  exit 1
fi

# Set username
USERNAME="dp-cli"

# Check if user already exists
if id "$USERNAME" &>/dev/null; then
  echo "User '$USERNAME' already exists."
else
  # Create a system user with:
  # - no home directory (-M)
  # - no login shell (--shell /usr/sbin/nologin)
  # - system account (--system)
  useradd --system --no-create-home --shell /usr/sbin/nologin "$USERNAME"
  echo "System user '$USERNAME' created successfully."
fi
