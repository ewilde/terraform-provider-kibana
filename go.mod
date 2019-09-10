module github.com/ewilde/terraform-provider-kibana

go 1.12

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/elazarl/goproxy v0.0.0-20190711103511-473e67f1d7d2 // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190711103511-473e67f1d7d2 // indirect
	github.com/ewilde/go-kibana v0.0.0-20190904184914-4cdb9115bcd1
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hashicorp/terraform v0.12.8
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2 // indirect
	github.com/ory/dockertest v3.3.5+incompatible // indirect
	github.com/pkg/errors v0.8.1
	google.golang.org/grpc v1.22.0 // indirect
)

replace github.com/ewilde/go-kibana => github.com/jfroche/go-kibana v0.0.0-20190902193413-dc51a6c88753
