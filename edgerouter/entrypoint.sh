#!/bin/sh

# Initialize the command with the executable
CMD="./edgerouter -i -d"

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
  if [ -s "$CLIENT_CERT_PATH" ]; then
    CMD="$CMD --client-cert $CLIENT_CERT_PATH"
  fi
fi

# Check if CLIENT_CERT_KEY_PATH is set and non-empty
if [ ! -z "$CLIENT_CERT_KEY_PATH" ]; then
  if [ -s "$CLIENT_CERT_KEY_PATH" ]; then
    CMD="$CMD --client-cert-key $CLIENT_CERT_KEY_PATH"
  fi
fi

# Check if CERT_PATH is set and non-empty
if [ ! -z "$CERT_PATH" ]; then
  if [ -s "$CERT_PATH" ]; then
    CMD="$CMD --cert-path $CERT_PATH"
  fi
fi

# Check if CERT_KEY_PATH is set and non-empty
if [ ! -z "$CERT_KEY_PATH" ]; then
  if [ -s "$CERT_PATH" ]; then
    CMD="$CMD --cert-key-path $CERT_KEY_PATH"
  fi
fi

# Check if CLIENT_CA is set and non-empty
if [ ! -z "$CLIENT_CA" ]; then
  if [ -s "$CLIENT_CA" ]; then
    CMD="$CMD --client-ca $CLIENT_CA"
  FI
fi

# Execute the constructed command
exec $CMD