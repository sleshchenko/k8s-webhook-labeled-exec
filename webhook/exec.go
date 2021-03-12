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
package webhook

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var V1PodExecOptionKind = metav1.GroupVersionKind{Kind: "PodExecOptions", Group: "", Version: "v1"}

const (
	ValidateWebhookPath = "/validate"
)

type WebhookHandler struct {
	Client  client.Client
	Decoder *admission.Decoder
}

// ResourcesValidator validates execs process all exec requests and:
type ResourcesValidator struct {
	*WebhookHandler
}

func NewResourcesValidator() *ResourcesValidator {
	return &ResourcesValidator{&WebhookHandler{}}
}

func (v *ResourcesValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	if req.Kind == V1PodExecOptionKind && req.Operation == v1beta1.Connect {
		p := corev1.Pod{}
		err := v.Client.Get(ctx, types.NamespacedName{
			Name:      req.Name,
			Namespace: req.Namespace,
		}, &p)

		if err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}

		v, ok := p.Labels["exec-defender.sleshche.com"]
		if ok && v == "activated" {
			return admission.Denied("You can't connect to pods which are labeled with `exec-defender.sleshche.com: activated`")
		}

		return admission.Allowed("Pod is not marked to prevent exec")
	}
	// Do not allow operation if the corresponding handler is not found
	// It indicates that the webhooks configuration is not a valid or incompatible with this version of controller
	return admission.Denied(fmt.Sprintf("This admission controller is not designed to handle %s operation for %s. Notify an administrator about this issue", req.Operation, req.Kind))
}

// WorkspaceMutator implements inject.Client.
// A client will be automatically injected.

// InjectClient injects the client.
func (v *ResourcesValidator) InjectClient(c client.Client) error {
	v.Client = c
	return nil
}

// WorkspaceMutator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *ResourcesValidator) InjectDecoder(d *admission.Decoder) error {
	v.Decoder = d
	return nil
}
