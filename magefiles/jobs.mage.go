package main

import (
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/helm"
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/k8s"
	// "github.com/DelineaXPM/dsv-k8s/v2/magefiles/kind"
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/minikube"
	"github.com/magefile/mage/mg"
	"github.com/pterm/pterm"
)

// Job is a namespace to contain chained sets of automation actions, to reduce the need to chain many commands together for common workflows.
type Job mg.Namespace

// Init runs the setup tasks to initialize the local resources and files, without trying to apply yet.
//
// Setup initializes all the required steps for the cluster creation, initial helm chart copies, and kubeconfig copies.
func (Job) Init() {
	pterm.DefaultSection.Println("(Job) Init()")
	mg.SerialDeps(
		// kind.Kind{}.Init,
		minikube.Minikube{}.Init,
		k8s.K8s{}.Init,
		helm.Helm{}.Init,
	)
}

// Setup runs all the initialization tasks, in addition to installing the helm chart and kubeconfig files.
func (Job) Setup() {
	pterm.DefaultSection.Println("(Job) Setup()")
	mg.SerialDeps(
		minikube.Minikube{}.Init,
		k8s.K8s{}.Init,
		helm.Helm{}.Init,
		mg.F(k8s.K8s{}.Apply, constants.CacheManifestDirectory),
		helm.Helm{}.Install,
		k8s.K8s{}.Logs,
	)
}

// Redeploy removes kubernetes resources and helm charts and then redeploys with log streaming by default.
func (Job) Redeploy() {
	pterm.DefaultSection.Println("(Job) Redeploy()")
	mg.SerialDeps(
		helm.Helm{}.Uninstall,
		mg.F(k8s.K8s{}.Delete, constants.CacheManifestDirectory),
		mg.F(k8s.K8s{}.Apply, constants.CacheManifestDirectory),
		helm.Helm{}.Install,
		k8s.K8s{}.Logs,
	)
}
