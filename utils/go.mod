module github.com/maritimeconnectivity/MMS/utils

go 1.22.2

require (
	github.com/google/uuid v1.6.0
	github.com/maritimeconnectivity/MMS/mmtp v0.0.0
	nhooyr.io/websocket v1.8.11
	golang.org/x/crypto v0.22.0
	google.golang.org/protobuf v1.33.0
)


replace github.com/maritimeconnectivity/MMS/mmtp => ../mmtp
