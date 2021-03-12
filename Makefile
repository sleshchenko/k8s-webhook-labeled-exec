# Copyright (c) 2019-2021 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation
#

SHELL := bash
.SHELLFLAGS = -ec
.ONESHELL:

export IMG ?= quay.io/sleshche/podexec-defender

docker:
	docker build . -t ${IMG} -f build/Dockerfile
	docker push ${IMG}

install:
	kubectl create namespace podexec-defender || true
	kubectl apply -f ./deploy/

uninstall:
	kubectl delete namespace podexec-defender
	kubectl delete validatingwebhookconfiguration podexec-defender.sleshche.com
	kubectl delete clusterrole podexec-defender
	kubectl delete clusterrolebinding podexec-defender