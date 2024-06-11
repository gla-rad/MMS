#!/bin/sh

# Initialize the command with the executable
CMD="./edgerouter"

# Check if MRN is set and non-empty
if [ ! -z "$MRN" ]; then
  CMD="$CMD --mrn $MRN"
fi

# Check if RADDR is set and non-empty
if [ ! -z "$RADDR" ]; then
  CMD="$CMD --raddr $RADDR"
fi

# Check if PORT is set and non-empty
if [ ! -z "$PORT" ]; then
  CMD="$CMD --port $PORT"
fi

# Check if CLIENT_CERT_PATH is set and non-empty
if [ ! -z "$CLIENT_CERT_PATH" ]; then
  CMD="$CMD --client-cert $CLIENT_CERT_PATH"
fi

# Check if CLIENT_CERT_KEY_PATH is set and non-empty
if [ ! -z "$CLIENT_CERT_KEY_PATH" ]; then
  CMD="$CMD --client-cert-key $CLIENT_CERT_KEY_PATH"
fi

# Check if CERT_PATH is set and non-empty
if [ ! -z "$CERT_PATH" ]; then
  CMD="$CMD --cert-path $CERT_PATH"
fi

# Check if CERT_KEY_PATH is set and non-empty
if [ ! -z "$CERT_KEY_PATH" ]; then
  CMD="$CMD --cert-key-path $CERT_KEY_PATH"
fi

# Check if CLIENT_CA is set and non-empty
if [ ! -z "$CLIENT_CA" ]; then
  CMD="$CMD --client-ca $CLIENT_CA"
fi

# Execute the constructed command
exec $CMD