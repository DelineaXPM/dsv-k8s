module github.com/DelineaXPM/dsv-k8s/v2

require (
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/DelineaXPM/dsv-sdk-go/v2 v2.0.0
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
)

replace dsv-k8s/internal/k8s => ./internal/k8s

replace dsv-k8s/pkg/config => ./pkg/config

replace dsv-k8s/pkg/injector => ./pkg/injector

replace dsv-k8s/pkg/patch => ./pkg/patch

replace dsv-k8s/pkg/syncer => ./pkg/syncer

go 1.16
