#!/bin/bash

# Check if root
if [ "$(id -u)" -ne "0" ]; then
  echo "Please run this script as root or using sudo."
  exit 1
fi

# Check and stop service on port 53
if netstat -tuln | grep ":53 " > /dev/null; then
  echo "Stopping service on port 53"
  service=$(netstat -tuln | grep ":53 " | awk '{print $7}')
  pid=$(echo $service | cut -d'/' -f1)
  kill -9 $pid
  echo "Service on port 53 stopped."
else
  echo "No service found on port 53"
fi
