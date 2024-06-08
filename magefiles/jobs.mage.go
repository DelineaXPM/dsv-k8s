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

// Redeploy removes k8s resources, helm uninstall, and then runs k8s apply and helm install.
func (Job) Redeploy() {
	pterm.DefaultSection.Println("(Job) Redeploy()")

	mg.Deps(
		helm.Helm{}.Uninstall,
		mg.F(k8s.K8s{}.Delete, constants.CacheManifestDirectory),
	)

	mg.SerialDeps(
		minikube.Minikube{}.RemoveImages,
		minikube.Minikube{}.LoadImages, // just be sure in case forget to load local images that the latest is always used
		helm.Helm{}.Install,            // this should take place first so the creation of the manifests can benefit from the resulting injector/syncer
		mg.F(k8s.K8s{}.Apply, constants.CacheManifestDirectory),
		// k8s.K8s{}.Logs, // use chained command
	)
}

// RebuildImages runs the build and minikube load commands so the new source is able to be run by `job:redeploy`.
func (Job) RebuildImages() {
	pterm.DefaultSection.Println("(Job) RebuildImages()")
	mg.SerialDeps(
		BuildAll,
		minikube.Minikube{}.LoadImages,
	)
	pterm.Success.Printfln("RebuildImages() complete. Run `mage job:redeploy` to redeploy the new images.")
}
