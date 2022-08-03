module github.com/DelineaXPM/dsv-k8s/v2

require (
	github.com/DelineaXPM/dsv-sdk-go/v2 v2.0.0
	github.com/magefile/mage v1.13.0
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/mittwald/go-helm-client v0.11.3
	github.com/pterm/pterm v0.12.45
	github.com/sheldonhull/magetools v0.0.10
	k8s.io/api v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v1.5.2
)

replace dsv-k8s/internal/k8s => ./internal/k8s

replace dsv-k8s/pkg/config => ./pkg/config

replace dsv-k8s/pkg/injector => ./pkg/injector

replace dsv-k8s/pkg/patch => ./pkg/patch

replace dsv-k8s/pkg/syncer => ./pkg/syncer

go 1.16
