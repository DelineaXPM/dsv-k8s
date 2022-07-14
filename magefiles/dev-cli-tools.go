package main

// tools is a list of tools that are installed as binaries for development usage.
// This list gets installed to go bin directory once `mage init` is run.
// This is for binaries that need to be invoked as cli tools, not packages.
var toolList = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup
	"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	"mvdan.cc/gofumpt@latest",
	"github.com/daixiang0/gci@latest",
	"github.com/goreleaser/goreleaser@latest",
	"github.com/iwittkau/mage-select@latest",
	"github.com/mfridman/tparse@latest", // Tparse provides nice formatted go test console output.
	"github.com/rakyll/gotest@latest",   // Gotest is a wrapper for running Go tests via command line with support for colors to make it more readable.
	"gotest.tools/gotestsum@latest",     // Gotestsum provides improved console output for tests as well as additional test output for CI systems.

	"sigs.k8s.io/kind@latest",                     // Kind is a k8 development environment tool.
	"github.com/gechr/yamlfmt@latest",             // Yamlfmt provides formatting standards for yaml files.
	"github.com/c-bata/kube-prompt@latest",        // Kube-prompt is a k8s interactive tool using go-prompt library to allow better kubectl completion discovery.
	"github.com/norwoodj/helm-docs/cmd/helm-docs", // Helm-docs is a tool for generating helm charts documentation.
}
