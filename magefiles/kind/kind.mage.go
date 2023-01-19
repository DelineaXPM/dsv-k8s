// Kind package contains all the tasks for automation of kind cluster creation and tear down, and the required kubectl commands to correctly use this.
package kind

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// Kind contains the kind cli commands.
type Kind mg.Namespace

func createCluster() error {
	mtu.CheckPtermDebug()
	kindargs := []string{
		"create",
		"cluster",
		"--name", constants.KindClusterName,
		"--wait",
		"300s",
	}
	if os.Getenv("KIND_SETUP_CONFIG") != "" {
		pterm.Info.Printfln("KIND_SETUP_CONFIG: %s", os.Getenv("KIND_SETUP_CONFIG"))
		kindargs = append(kindargs, "--config", os.Getenv("KIND_SETUP_CONFIG"))
	}
	if err := sh.RunV(
		"kind",
		kindargs...,
	); err != nil {
		return err
	}
	return nil
}

func updateKubeconfig() error {
	mtu.CheckPtermDebug()
	if _, err := os.Stat(constants.Kubeconfig); os.IsNotExist(err) {
		if _, err := os.Create(constants.Kubeconfig); err != nil {
			pterm.Error.Printfln("unable to create empty placeholder file: %v", err)
		}
	}
	kc, err := sh.Output("kind", "get", "cluster", "kubeconfig", "--name", constants.KindClusterName)
	if err != nil {
		pterm.Error.Println("unable to get kind cluster info, maybe you need to run mage kind:init first?")
		return err
	}

	if err := os.WriteFile(constants.Kubeconfig, []byte(kc), constants.PermissionUserReadWriteExecute); err != nil {
		pterm.Error.Printfln("unable to write kubeconfig to file: %v", err)
		return err
	}
	pterm.Info.Printfln("kubeconfig updated: %s", constants.Kubeconfig)
	// for now this is only going to be run against Kind cluster.
	// if err := sh.Run(
	// 	"kubectl",
	// 	"cluster-info", "--context", constants.KindContextName,
	// 	"--cluster", constants.KindContextName,
	// ); err != nil {
	// 	return err
	// }
	return nil
}

// ➕ Create creates a new Kind cluster and populates a kubeconfig in cachedirectory.
func (Kind) Init() error {
	mtu.CheckPtermDebug()

	out, err := sh.Output("kind", "get", "clusters")
	if err := err; err != nil {
		return err
	}
	cleanOutput := strings.TrimSpace(out)
	matchedCluster := regexp.MustCompile(constants.KindClusterName)
	pterm.Debug.Printfln("cleanOutput: %s", cleanOutput)
	if !matchedCluster.MatchString(cleanOutput) {
		pterm.Info.Printfln("simple match not found, so attempting to recreate cluster")
		if err := createCluster(); err != nil {
			return err
		}
	}
	dspin, _ := pterm.DefaultSpinner.
		WithDelay(time.Second).
		WithRemoveWhenDone(true).
		WithShowTimer(true).
		WithText("Init()\n").
		WithSequence("|", "/", "-", "|", "/", "-", "\\").Start()
	dspin.SuccessPrinter.Println("ensuring it's in kubeconfig")
	if err := updateKubeconfig(); err != nil {
		pterm.Error.Printfln("updateKubeconfig(): %v", err)
	}
	dspin.UpdateText("setting context")
	if err := sh.Run("kubectl", "config", "use-context", constants.KindContextName); err != nil {
		dspin.WarningPrinter.Printfln("default context might not be setup correct to new context: %v", err)
	}
	if err := sh.Run("kubectl", "config", "set-context", "--context", constants.KindContextName, "--current", "--namespace", constants.KubectlNamespace); err != nil {
		dspin.WarningPrinter.Printfln("default namespace might not be setup correct to new namespace: %v", err)
	}
	// Create the namespace if it doesn't exist.
	dspin.UpdateText("creating namespace if not exists")
	if _, err := sh.Output("kubectl", "get", "namespace", constants.KubectlNamespace); err != nil {
		dspin.UpdateText(fmt.Sprintf("namespace does not exist, creating namespace: %s...", constants.KubectlNamespace))

		if err := sh.Run("kubectl", "create", "namespace", constants.KubectlNamespace); err != nil {
			dspin.FailPrinter.Printfln("unable to create namespace: %v", err)
			return fmt.Errorf("kubectl create namespace %s: %w", constants.KubectlNamespace, err)
		}
		dspin.SuccessPrinter.Printfln("namespace created: %s", constants.KubectlNamespace)
	}
	dspin.UpdateText("pulling docker images")
	if err := sh.Run("docker", "pull", "quay.io/delinea/dsv-k8s:latest"); err != nil {
		dspin.WarningPrinter.Printfln("docker pull: %v", err)
		return fmt.Errorf("docker pull: %w", err)
	}
	dspin.SuccessPrinter.Println("docker image for quay.io/delinea/dsv-k8s:latest pulled")
	// Not working right now, can't find nodes for Kind to preload. Not critical so commenting out for now - sheldon.
	// Sp.UpdateText("preloading docker image into kind cluster")
	// if err := sh.Run("kind", "load", "docker-image", "quay.io/delinea/dsv-k8s:latest"); err != nil {
	// 	return fmt.Errorf("kind load docker-image: %w", err)
	// }.
	dspin.SuccessPrinter.Println("(Kind) Init()")
	_ = dspin.Stop()
	return nil
}

// 🗑️ Destroy tears down the Kind cluster.
func (Kind) Destroy() error {
	mtu.CheckPtermDebug()
	if err := sh.Run("kind", "delete", "cluster", "--name", constants.KindClusterName); err != nil {
		pterm.Error.Printfln("kind delete error: %v", err)
		return err
	}
	if err := sh.Run("kubectl", "config", "unset", fmt.Sprintf("clusters.%s", constants.KindContextName)); err != nil {
		pterm.Warning.Printfln("default context might not be setup correct to new context: %v", err)
	}

	pterm.Success.Println("(Kind) Destroy()")
	return nil
}
