module github.com/maritimeconnectivity/MMS/edgerouter

go 1.21

require (
	github.com/google/uuid v1.3.1
	github.com/hashicorp/mdns v1.0.5
	github.com/maritimeconnectivity/MMS/mmtp v0.0.0
	golang.org/x/crypto v0.13.0
	google.golang.org/protobuf v1.31.0
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/miekg/dns v1.1.41 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

replace github.com/maritimeconnectivity/MMS/mmtp => ../mmtp
