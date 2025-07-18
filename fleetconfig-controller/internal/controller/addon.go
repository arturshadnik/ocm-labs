package controller

import (
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	corev1 "k8s.io/api/core/v1"
	kerrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	addonapi "open-cluster-management.io/api/client/addon/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/open-cluster-management-io/lab/fleetconfig-controller/api/v1alpha1"
	exec_utils "github.com/open-cluster-management-io/lab/fleetconfig-controller/internal/exec"
	"github.com/open-cluster-management-io/lab/fleetconfig-controller/internal/file"
)

const (
	addonConfigMapName = "fleet-add-ons"
	addon              = "addon"
	create             = "create"
	enable             = "enable"
	disable            = "disable"
)

type manifestType string

const (
	manifestTypeRaw  manifestType = "raw"
	manifestTypeHTTP manifestType = "http"
)

func handleAddonConfig(ctx context.Context, kClient client.Client, addonC *addonapi.Clientset, fc *v1alpha1.FleetConfig) error {
	logger := log.FromContext(ctx)
	logger.V(0).Info("handleAddOnConfig", "fleetconfig", fc.Name)

	// get existing addons
	createdAddOns, err := addonC.AddonV1alpha1().AddOnTemplates().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	requestedAddOns := fc.Spec.AddOnConfigs

	// nothing to do
	if len(requestedAddOns) == 0 && len(createdAddOns.Items) == 0 {
		logger.V(5).Info("no addons to reconcile")
		return nil
	}

	// compare existing to requested
	createdVersionedNames := make([]string, len(createdAddOns.Items))
	for i, ca := range createdAddOns.Items {
		createdVersionedNames[i] = ca.Name
	}

	requestedVersionedNames := make([]string, len(requestedAddOns))
	for i, ra := range requestedAddOns {
		requestedVersionedNames[i] = fmt.Sprintf("%s-%s", ra.Name, ra.Version)
	}

	// Find addons that need to be created (present in requested, missing from created)
	addonsToCreate := make([]*v1alpha1.AddOnConfig, 0)
	for i, requestedName := range requestedVersionedNames {
		if !slices.Contains(createdVersionedNames, requestedName) {
			addonsToCreate = append(addonsToCreate, requestedAddOns[i])
		}
	}

	// Find addons that need to be deleted (present in created, missing from requested)
	addonsToDelete := make([]string, 0)
	for _, createdName := range createdVersionedNames {
		if !slices.Contains(requestedVersionedNames, createdName) {
			addonsToDelete = append(addonsToDelete, createdName)
		}
	}

	// do deletes first, then creates.
	err = handleAddonDelete(ctx, addonC, fc, addonsToDelete)
	if err != nil {
		return err
	}

	err = handleAddonCreate(ctx, kClient, fc, addonsToCreate)
	if err != nil {
		return err
	}

	return nil
}

func handleAddonCreate(ctx context.Context, kClient client.Client, fc *v1alpha1.FleetConfig, addons []*v1alpha1.AddOnConfig) error {
	if len(addons) == 0 {
		return nil
	}

	logger := log.FromContext(ctx)
	logger.V(0).Info("createAddOns", "fleetconfig", fc.Name)

	// look up CM
	addonConfigMap := corev1.ConfigMap{}
	err := kClient.Get(ctx, types.NamespacedName{Name: addonConfigMapName, Namespace: fc.Namespace}, &addonConfigMap)
	if err != nil {
		return err
	}

	data := addonConfigMap.Data

	// set up array of clusteradm addon create commands
	for _, a := range addons {
		// pull out manifests
		manifests, ok := data[fmt.Sprintf("%s-%s", a.Name, a.Version)]
		if !ok {
			return fmt.Errorf("no manifests found for add-on %s version %s", a.Name, a.Version)
		}

		args := []string{
			addon,
			create,
			a.Name,
			fmt.Sprintf("--version=%s", a.Version),
		}

		switch getManifestType(manifests) {
		case manifestTypeHTTP:
			// pass URL directly
			args = append(args, fmt.Sprintf("--filename=%s", manifests))
		case manifestTypeRaw:
			filename, cleanup, err := file.TmpFile([]byte(manifests), "yaml")
			if cleanup != nil {
				defer cleanup()
			}
			if err != nil {
				return err
			}
			args = append(args, fmt.Sprintf("--filename=%s", filename))
		}

		if a.HubRegistration {
			args = append(args, "--hub-registration")
		}
		if a.Overwrite {
			args = append(args, "--overwrite")
		}
		if a.ClusterRoleBinding != "" {
			args = append(args, fmt.Sprintf("--cluster-role-bind=%s", a.ClusterRoleBinding))
		}

		cmd := exec.Command(clusteradm, args...)
		out, err := exec_utils.CmdWithLogs(ctx, cmd, "waiting for 'clusteradm addon create' to complete...")
		if err != nil {
			return fmt.Errorf("failed to create addon: %v, output: %s", err, string(out))
		}
		logger.V(0).Info("created addon", "AddOnTemplate", a.Name)
	}
	return nil
}

func getManifestType(manifests string) manifestType {
	if strings.HasPrefix(manifests, "http://") || strings.HasPrefix(manifests, "https://") {
		return manifestTypeHTTP
	}
	return manifestTypeRaw
}

func handleAddonDelete(ctx context.Context, addonC *addonapi.Clientset, fc *v1alpha1.FleetConfig, addons []string) error {
	if len(addons) == 0 {
		return nil
	}

	logger := log.FromContext(ctx)
	logger.V(0).Info("deleteAddOns", "fleetconfig", fc.Name)

	// a list of addons which may or may not need to be purged at the end (ClusterManagementAddOns needs to be deleted)
	purgeList := make([]string, 0)
	for _, addonName := range addons {
		// get the addon template, so we can extract spec.addonName
		addon, err := addonC.AddonV1alpha1().AddOnTemplates().Get(ctx, addonName, metav1.GetOptions{})
		if err != nil && !kerrs.IsNotFound(err) {
			return fmt.Errorf("failed to delete addon %s: %v", addonName, err)
		}

		// delete the addon template
		if addon != nil {
			err = addonC.AddonV1alpha1().AddOnTemplates().Delete(ctx, addonName, metav1.DeleteOptions{})
			if err != nil && !kerrs.IsNotFound(err) {
				return fmt.Errorf("failed to delete addon %s: %v", addonName, err)
			}
		}

		// get the addon name without a version suffix, add it to purge list
		purgeList = append(purgeList, addon.Spec.AddonName)
		logger.V(0).Info("deleted addon", "AddOnTemplate", addonName)
	}

	// check if there are any remaining addon templates for the same addon names as what was just deleted (different versions of the same addon)
	allAddons, err := addonC.AddonV1alpha1().AddOnTemplates().List(ctx, metav1.ListOptions{})
	if err != nil && !kerrs.IsNotFound(err) {
		return fmt.Errorf("failed to clean up addons %v: %v", purgeList, err)
	}
	for _, a := range allAddons.Items {
		// if other versions of the same addon remain, remove it from the purge list
		purgeList = slices.DeleteFunc(purgeList, func(name string) bool {
			return name == a.Spec.AddonName
		})
	}
	// if list is empty, nothing else to do
	if len(purgeList) == 0 {
		return nil
	}

	// delete the ClusterManagementAddOn for any addon which has no active versions left
	for _, name := range purgeList {
		err = addonC.AddonV1alpha1().ClusterManagementAddOns().Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil && !kerrs.IsNotFound(err) {
			return fmt.Errorf("failed to purge addon %s: %v", name, err)
		}
		logger.V(0).Info("purged addon", "ClusterManagementAddOn", name)
	}

	return nil
}
