#!/bin/bash

# Must be run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root." >&2
  exit 1
fi

# Set the username to delete
USERNAME="dp-cli"

# Check if the user exists
if id "$USERNAME" &>/dev/null; then
  # Delete the user and their home directory if it existed
  userdel "$USERNAME"
  echo "System user '$USERNAME' has been deleted."
else
  echo "User '$USERNAME' does not exist."
fi
