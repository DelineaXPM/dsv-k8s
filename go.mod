module github.com/thycotic/dsv-k8s

require (
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/thycotic/dsv-sdk-go v1.0.1
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
)

// replace github.com/thycotic/dsv-sdk-go => ../dsv-sdk-go
replace dsv-k8s/internal/k8s => ./internal/k8s

replace dsv-k8s/pkg/config => ./pkg/config

replace dsv-k8s/pkg/injector => ./pkg/injector

replace dsv-k8s/pkg/patch => ./pkg/patch

replace dsv-k8s/pkg/syncer => ./pkg/syncer

go 1.16
