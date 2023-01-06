package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Cluster"
// +kubebuilder:printcolumn:name="ADDON NAME",type=string,JSONPath=`.spec.addonName`

// AddOnTemplate is the Custom Resource object, it is used to describe
// how to deploy the addon agent and how to register the addon.
//
// AddOnTemplate is a cluster-scoped resource, and will only be used
// on the hub cluster.
type AddOnTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// spec holds the registration configuration for the addon and the
	// addon agent resources yaml description.
	// +kubebuilder:validation:Required
	// +required
	Spec AddOnTemplateSpec `json:"spec"`
}

// AddOnTemplateSpec defines the template of an addon agent which will be deployed on managed clusters.
type AddOnTemplateSpec struct {
	// AddonName represents the name of the addon which the template belongs to
	// +kubebuilder:validation:Required
	// +required
	AddonName string `json:"addonName"`

	// AgentManifests represents the kubernetes resources of the addon agent to be deployed on a managed cluster.
	// +kubebuilder:validation:Required
	// +required
	AgentManifests []Manifest `json:"agentManifests"`

	// Registration holds the registration configuration for the addon
	// +kubebuilder:validation:Required
	// +required
	Registration []RegistrationSpec `json:"registration"`
}

// Manifest represents a resource to be deployed on the managed cluster.
type Manifest struct {
	// +kubebuilder:validation:EmbeddedResource
	// +kubebuilder:pruning:PreserveUnknownFields
	runtime.RawExtension `json:",inline"`
}

// RegistrationType represents the type of the registration configuration,
// it could be KubeClient or CustomSigner
type RegistrationType string

// CSRApproveStrategyType represent how to approve the addon registration
// Certificate Signing Requests
type CSRApproveStrategyType string

const (
	// RegistrationTypeKubeClient represents the KubeClient type registration of the addon agent.
	// For this type, the addon agent can access the hub kube apiserver with kube style API.
	// The signer name should be "kubernetes.io/kube-apiserver-client".
	RegistrationTypeKubeClient RegistrationType = "KubeClient"
	// RegistrationTypeCustomSigner represents the CustomSigner type registration of the addon agent.
	// For this type, the addon agent can access the hub cluster through user-defined endpoints.
	RegistrationTypeCustomSigner RegistrationType = "CustomSigner"

	// CSRApproveStrategyAuto means automatically approve the CSR
	CSRApproveStrategyAuto CSRApproveStrategyType = "Auto"
	// CSRApproveStrategyNone means that the CSR will not be approved
	// automatically, users need to approve them by themselves
	CSRApproveStrategyNone CSRApproveStrategyType = "None"
)

// RegistrationSpec describes how to register an addon agent to the hub cluster.
// With the registration defined, The addon agent can access to kube apiserver with kube style API
// or other endpoints on hub cluster with client certificate authentication. During the addon
// registration process, a csr will be created for each RegistrationSpec on the hub cluster. The
// CSR can be approved automatically(Auto) or manually(None), After the csr is approved on the hub
// cluster, the klusterlet agent will create a secret in the installNamespace for the addon agent.
// If the RegistrationType type is KubeClient, the secret name will be "{addon name}-hub-kubeconfig"
// whose content includes key/cert and kubeconfig. Otherwise, If the RegistrationType type is
// CustomSigner the secret name will be "{addon name}-{signer name}-client-cert" whose content
// includes key/cert.
type RegistrationSpec struct {
	// Type of the registration configuration
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=KubeClient;CustomSigner
	Type RegistrationType `json:"type"`

	// ApproveStrategy represents how to approve the addon registration.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=Auto;None
	ApproveStrategy CSRApproveStrategyType `json:"approveStrategy"`

	// KubeClient holds the configuration of the KubeClient type registration
	// +optional
	KubeClient *KubeClientRegistrationConfig `json:"kubeClient,omitempty"`

	// KubeClient holds the configuration of the CustomSigner type registration
	// required when the Type is CustomSigner
	CustomSigner *CustomSignerRegistrationConfig `json:"customSigner,omitempty"`
}

type KubeClientRegistrationConfig struct {
	// Permission represents the permission configuration of the addon agent to access the hub cluster
	// +optional
	Permission *HubPermissionConfig `json:"permission,omitempty"`
}

// HubPermissionConfig configures the permission of the addon agent to access the hub cluster.
// Will create a RoleBinding in the same namespace as the managedClusterAddon to bind the user
// provided ClusterRole/Role to the "system:open-cluster-management:cluster:<cluster-name>:addon:<addon-name>"
// Group.
type HubPermissionConfig struct {
	// ClusterRoleName of the permission setting cluster role.
	// +optional
	ClusterRoleName string `json:"clusterRoleName,omitempty"`
	// RoleName of the permission setting role in the same namespace as the managedClusterAddon.
	// +optional
	RoleName string `json:"roleName,omitempty"`
}

type CustomSignerRegistrationConfig struct {
	// Name of the signer
	// +required
	// +kubebuilder:validation:MaxLength=571
	// +kubebuilder:validation:MinLength=5
	Name string `json:"name"`
	// SigningCARef represents the reference of the secret to sign the CSR
	// +kubebuilder:validation:Required
	SigningCA SigningCARef `json:"signingCA"`
}

// SigningCARef is the reference to the signing CA secret
type SigningCARef struct {
	// Namespace of the signing CA secret
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
	// Name of the signing CA secret
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}
