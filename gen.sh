#!/usr/bin/env bash

if [ -z "$1" ]
then
  echo "Please supply a subdomain to create a certificate for";
  echo "e.g. www.mysite.com"
  exit -1;
fi

COMMON_NAME=$1
SUBJECT="/C=CA/ST=None/L=NB/O=None/CN=$COMMON_NAME"
NUM_OF_DAYS=825

cat v3.ext | sed s/%%DOMAIN%%/"$COMMON_NAME"/g > /tmp/__v3.ext
openssl req -new -newkey rsa:2048 -sha256 -nodes -key cert.key -config /tmp/__v3.ext -subj "$SUBJECT" -out device.csr
openssl x509 -req -in device.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days $NUM_OF_DAYS -sha256 -extfile /tmp/__v3.ext 

# remove temp file
rm -f device.csr
