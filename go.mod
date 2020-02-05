module github.com/thycotic/dsv-k8s

require (
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/thycotic/dsv-sdk-go v0.0.0-20200118074348-81e52d764b28
	k8s.io/api v0.0.0-20190720062849-3043179095b6
	k8s.io/apimachinery v0.0.0-20190719140911-bfcf53abc9f8
)

// replace github.com/thycotic/dsv-sdk-go => ../dsv-sdk-go
replace dsv-k8s/pkg/vault => ./pkg/vault

go 1.13
