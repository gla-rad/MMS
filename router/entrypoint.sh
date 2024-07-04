#!/bin/sh

# Initialize the command with the executable
CMD="./router"

# Check if PORT is set and non-empty
if [ ! -z "$PORT" ]; then
  CMD="$CMD --port $PORT"
fi

# Check if LIBP2P_PORT is set and non-empty
if [ ! -z "$LIBP2P_PORT" ]; then
  CMD="$CMD --libp2p-port $LIBP2P_PORT"
fi

# Check if PRIVKEY is set and non-empty
if [ ! -z "$PRIVKEY" ]; then
  CMD="$CMD --privkey $PRIVKEY"
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

# Check if BEACONS is set and non-empty
if [ ! -z "$BEACONS" ]; then
  CMD="$CMD --beacons $BEACONS"
fi

# Execute the constructed command
exec $CMD