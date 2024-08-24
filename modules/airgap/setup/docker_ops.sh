#!/bin/bash

set -x
echo "$@"

# Define the input file
infile=`pwd`/$1
hostdns=${2}
username=${3}
password=${4}

# Read the input file line by line using a for loop
IFS=$'\n' # set the Internal Field Separator to newline
for line in $(cat "$infile"); do
  if [[ "$line" =~ "docker" ]]; then
    line=`echo "${line/docker.io\/}"`
  fi
  docker pull $line && \
  docker image tag $line $hostdns/$line && \
  echo "$password" | docker login $hostdns -u "$username" --password-stdin && \
  docker push $hostdns/$line
  echo "Docker pull/tag/push completed for image: $line"
done