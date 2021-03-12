//
// Copyright (c) 2019-2021 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package cluster

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	WatchNamespaceEnvVar = "WATCH_NAMESPACE"
)

// GetWatchNamespace returns the namespace the operator should be watching for changes
//
// This function was ported over from Operator SDK 0.17.0
func GetWatchNamespace() (string, error) {
	ns, found := os.LookupEnv(WatchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", WatchNamespaceEnvVar)
	}
	return ns, nil
}

func findAPIResources(source []*metav1.APIResourceList, groupName string) []metav1.APIResource {
	for i := 0; i < len(source); i++ {
		if source[i].GroupVersion == groupName {
			return source[i].APIResources
		}
	}
	return nil
}

//IsWebhookConfigurationEnabled returns true if both of mutating and validating webhook configurations are enabled
func IsWebhookConfigurationEnabled() (bool, error) {
	kubeCfg, err := config.GetConfig()
	if err != nil {
		return false, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeCfg)
	if err != nil {
		return false, err
	}
	_, apiResources, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		return false, err
	}

	if admissionRegistrationResources := findAPIResources(apiResources, "admissionregistration.k8s.io/v1beta1"); admissionRegistrationResources != nil {
		isMutatingHookAvailable := false
		isValidatingMutatingHookAvailable := false
		for i := range admissionRegistrationResources {
			if admissionRegistrationResources[i].Name == "mutatingwebhookconfigurations" {
				isMutatingHookAvailable = true
			}

			if admissionRegistrationResources[i].Name == "validatingwebhookconfigurations" {
				isValidatingMutatingHookAvailable = true
			}
		}

		return isMutatingHookAvailable && isValidatingMutatingHookAvailable, nil
	}

	return false, nil
}
