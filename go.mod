module github.com/ewilde/terraform-provider-kibana

go 1.13

require (
	github.com/Microsoft/go-winio v0.4.15 // indirect
	github.com/apparentlymart/go-dump v0.0.0-20190214190832-042adf3cf4a0 // indirect
	github.com/aws/aws-sdk-go v1.25.3 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20200928162600-f2cc35102c2a // indirect
	github.com/ewilde/go-kibana v0.0.0-20210127120218-80bc38c8b5b8
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/hashicorp/hcl v0.0.0-20171017181929-23c074d0eceb // indirect
	github.com/hashicorp/hcl2 v0.0.0-20190821123243-0c888d1241f6 // indirect
	github.com/hashicorp/hil v0.0.0-20190212112733-ab17b08d6590 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/opencontainers/runc v1.0.0-rc4.0.20171130145147-91e979501348 // indirect
	github.com/ory/dockertest v3.3.5+incompatible // indirect
	github.com/parnurzeal/gorequest v0.2.16 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/vmihailenco/msgpack v4.0.1+incompatible // indirect
	golang.org/x/crypto v0.0.0-20201117144127-c1f2f97bffc9 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	moul.io/http2curl v1.0.0 // indirect

)

replace (
	github.com/ewilde/go-kibana => /home/peter/tools/go-kibana
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190830141801-acfa387b8d69
)
