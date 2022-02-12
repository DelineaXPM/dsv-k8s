module github.com/thycotic/dsv-k8s

require (
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/thycotic/dsv-sdk-go v1.0.1
	k8s.io/api v0.23.3
	k8s.io/apimachinery v0.23.3
)

// replace github.com/thycotic/dsv-sdk-go => ../dsv-sdk-go
replace dsv-k8s/pkg/vault => ./pkg/vault

go 1.16
