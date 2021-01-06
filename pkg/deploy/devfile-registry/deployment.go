//
// Copyright (c) 2020-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//
package devfile_registry

import (
	"github.com/eclipse/che-operator/pkg/deploy"
	"github.com/eclipse/che-operator/pkg/deploy/registry"
	"github.com/eclipse/che-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
)

func SyncDevfileRegistryDeploymentToCluster(deployContext *deploy.DeployContext) deploy.DeploymentProvisioningStatus {
	registryType := "devfile"
	registryImage := util.GetValue(deployContext.CheCluster.Spec.Server.DevfileRegistryImage, deploy.DefaultDevfileRegistryImage(deployContext.CheCluster))
	registryImagePullPolicy := v1.PullPolicy(util.GetValue(string(deployContext.CheCluster.Spec.Server.DevfileRegistryPullPolicy), deploy.DefaultPullPolicyFromDockerImage(registryImage)))
	probePath := "/devfiles/"
	devfileImagesEnv := util.GetEnvByRegExp("^.*devfile_registry_image.*$")

	clusterDeployment, err := deploy.GetClusterDeployment(deploy.DevfileRegistry, deployContext.CheCluster.Namespace, deployContext.ClusterAPI.Client)
	if err != nil {
		return deploy.DeploymentProvisioningStatus{
			ProvisioningStatus: deploy.ProvisioningStatus{Err: err},
		}
	}

	resources := v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceMemory: util.GetResourceQuantity(
				deployContext.CheCluster.Spec.Server.DevfileRegistryMemoryRequest,
				deploy.DefaultDevfileRegistryMemoryRequest),
		},
		Limits: v1.ResourceList{
			v1.ResourceMemory: util.GetResourceQuantity(
				deployContext.CheCluster.Spec.Server.DevfileRegistryMemoryLimit,
				deploy.DefaultDevfileRegistryMemoryLimit),
			v1.ResourceCPU: util.GetResourceQuantity(
				deployContext.CheCluster.Spec.Server.DevfileRegistryCpuLimit,
				deploy.DefaultDevfileRegistryCpuLimit),
		},
	}

	specDeployment, err := registry.GetSpecRegistryDeployment(
		deployContext,
		registryType,
		registryImage,
		devfileImagesEnv,
		registryImagePullPolicy,
		resources,
		probePath)

	if err != nil {
		return deploy.DeploymentProvisioningStatus{
			ProvisioningStatus: deploy.ProvisioningStatus{Err: err},
		}
	}

	return deploy.SyncDeploymentToCluster(deployContext, specDeployment, clusterDeployment, nil, nil)
}
