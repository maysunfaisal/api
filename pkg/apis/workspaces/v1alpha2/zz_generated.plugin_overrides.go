package v1alpha2

import (
	attributes "github.com/devfile/api/v2/pkg/attributes"
)

// +devfile:jsonschema:generate
type PluginOverrides struct {
	OverridesBase `json:",inline"`

	// Overrides of attributes encapsulated in a parent devfile or a plugin.
	// Overriding is done according to K8S strategic merge patch standard rules.
	// +optional
	// +patchStrategy=merge
	Attributes attributes.Attributes `json:"attributes,omitempty" patchStrategy:"merge"`

	// Overrides of components encapsulated in a parent devfile or a plugin.
	// Overriding is done according to K8S strategic merge patch standard rules.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +devfile:toplevellist
	Components []ComponentPluginOverride `json:"components,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Overrides of commands encapsulated in a parent devfile or a plugin.
	// Overriding is done according to K8S strategic merge patch standard rules.
	// +optional
	// +patchMergeKey=id
	// +patchStrategy=merge
	// +devfile:toplevellist
	Commands []CommandPluginOverride `json:"commands,omitempty" patchStrategy:"merge" patchMergeKey:"id"`
}

//+k8s:openapi-gen=true
type ComponentPluginOverride struct {

	// Mandatory name that allows referencing the component
	// from other elements (such as commands) or from an external
	// devfile that may reference this component through a parent or a plugin.
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`

	// Map of implementation-dependant free-form YAML attributes.
	// +optional
	Attributes                   attributes.Attributes `json:"attributes,omitempty"`
	ComponentUnionPluginOverride `json:",inline"`
}

type CommandPluginOverride struct {

	// Mandatory identifier that allows referencing
	// this command in composite commands, from
	// a parent, or in events.
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// +kubebuilder:validation:MaxLength=63
	Id string `json:"id"`

	// Map of implementation-dependant free-form YAML attributes.
	// +optional
	Attributes                 attributes.Attributes `json:"attributes,omitempty"`
	CommandUnionPluginOverride `json:",inline"`
}

// +union
type ComponentUnionPluginOverride struct {

	// +kubebuilder:validation:Enum=Container;Kubernetes;Openshift;Volume
	// Type of component
	//
	// +unionDiscriminator
	// +optional
	ComponentType ComponentTypePluginOverride `json:"componentType,omitempty"`

	// Allows adding and configuring workspace-related containers
	// +optional
	Container *ContainerComponentPluginOverride `json:"container,omitempty"`

	// Allows importing into the workspace the Kubernetes resources
	// defined in a given manifest. For example this allows reusing the Kubernetes
	// definitions used to deploy some runtime components in production.
	//
	// +optional
	Kubernetes *KubernetesComponentPluginOverride `json:"kubernetes,omitempty"`

	// Allows importing into the workspace the OpenShift resources
	// defined in a given manifest. For example this allows reusing the OpenShift
	// definitions used to deploy some runtime components in production.
	//
	// +optional
	Openshift *OpenshiftComponentPluginOverride `json:"openshift,omitempty"`

	// Allows specifying the definition of a volume
	// shared by several other components
	// +optional
	Volume *VolumeComponentPluginOverride `json:"volume,omitempty"`
}

// +union
type CommandUnionPluginOverride struct {

	// +kubebuilder:validation:Enum=Exec;Apply;Composite
	// Type of workspace command
	// +unionDiscriminator
	// +optional
	CommandType CommandTypePluginOverride `json:"commandType,omitempty"`

	// CLI Command executed in an existing component container
	// +optional
	Exec *ExecCommandPluginOverride `json:"exec,omitempty"`

	// Command that consists in applying a given component definition,
	// typically bound to a workspace event.
	//
	// For example, when an `apply` command is bound to a `preStart` event,
	// and references a `container` component, it will start the container as a
	// K8S initContainer in the workspace POD, unless the component has its
	// `dedicatedPod` field set to `true`.
	//
	// When no `apply` command exist for a given component,
	// it is assumed the component will be applied at workspace start
	// by default.
	// +optional
	Apply *ApplyCommandPluginOverride `json:"apply,omitempty"`

	// Composite command that allows executing several sub-commands
	// either sequentially or concurrently
	// +optional
	Composite *CompositeCommandPluginOverride `json:"composite,omitempty"`
}

// ComponentType describes the type of component.
// Only one of the following component type may be specified.
type ComponentTypePluginOverride string

// Component that allows the developer to add a configured container into his workspace
type ContainerComponentPluginOverride struct {
	BaseComponentPluginOverride `json:",inline"`
	ContainerPluginOverride     `json:",inline"`
	Endpoints                   []EndpointPluginOverride `json:"endpoints,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// Component that allows partly importing Kubernetes resources into the workspace POD
type KubernetesComponentPluginOverride struct {
	K8sLikeComponentPluginOverride `json:",inline"`
}

// Component that allows partly importing Openshift resources into the workspace POD
type OpenshiftComponentPluginOverride struct {
	K8sLikeComponentPluginOverride `json:",inline"`
}

// Component that allows the developer to declare and configure a volume into his workspace
type VolumeComponentPluginOverride struct {
	BaseComponentPluginOverride `json:",inline"`
	VolumePluginOverride        `json:",inline"`
}

// CommandType describes the type of command.
// Only one of the following command type may be specified.
type CommandTypePluginOverride string

type ExecCommandPluginOverride struct {
	LabeledCommandPluginOverride `json:",inline"`

	//  +optional
	// The actual command-line string
	//
	// Special variables that can be used:
	//
	//  - `$PROJECTS_ROOT`: A path where projects sources are mounted as defined by container component's sourceMapping.
	//
	//  - `$PROJECT_SOURCE`: A path to a project source ($PROJECTS_ROOT/<project-name>). If there are multiple projects, this will point to the directory of the first one.
	CommandLine string `json:"commandLine,omitempty"`

	//  +optional
	// Describes component to which given action relates
	//
	Component string `json:"component,omitempty"`

	// Working directory where the command should be executed
	//
	// Special variables that can be used:
	//
	//  - `$PROJECTS_ROOT`: A path where projects sources are mounted as defined by container component's sourceMapping.
	//
	//  - `$PROJECT_SOURCE`: A path to a project source ($PROJECTS_ROOT/<project-name>). If there are multiple projects, this will point to the directory of the first one.
	// +optional
	WorkingDir string `json:"workingDir,omitempty"`

	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// Optional list of environment variables that have to be set
	// before running the command
	Env []EnvVarPluginOverride `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// +optional
	// Whether the command is capable to reload itself when source code changes.
	// If set to `true` the command won't be restarted and it is expected to handle file changes on its own.
	//
	// Default value is `false`
	HotReloadCapable bool `json:"hotReloadCapable,omitempty"`
}

type ApplyCommandPluginOverride struct {
	LabeledCommandPluginOverride `json:",inline"`

	//  +optional
	// Describes component that will be applied
	//
	Component string `json:"component,omitempty"`
}

type CompositeCommandPluginOverride struct {
	LabeledCommandPluginOverride `json:",inline"`

	// The commands that comprise this composite command
	Commands []string `json:"commands,omitempty" patchStrategy:"replace"`

	// Indicates if the sub-commands should be executed concurrently
	// +optional
	Parallel bool `json:"parallel,omitempty"`
}

// Workspace component: Anything that will bring additional features / tooling / behaviour / context
// to the workspace, in order to make working in it easier.
type BaseComponentPluginOverride struct {
}

type ContainerPluginOverride struct {
	//  +optional
	Image string `json:"image,omitempty"`

	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// Environment variables used in this container.
	//
	// The following variables are reserved and cannot be overridden via env:
	//
	//  - `$PROJECTS_ROOT`
	//
	//  - `$PROJECT_SOURCE`
	Env []EnvVarPluginOverride `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// +optional
	// List of volumes mounts that should be mounted is this container.
	VolumeMounts []VolumeMountPluginOverride `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// +optional
	MemoryLimit string `json:"memoryLimit,omitempty"`

	// +optional
	MemoryRequest string `json:"memoryRequest,omitempty"`

	// +optional
	CpuLimit string `json:"cpuLimit,omitempty"`

	// +optional
	CpuRequest string `json:"cpuRequest,omitempty"`

	// The command to run in the dockerimage component instead of the default one provided in the image.
	//
	// Defaults to an empty array, meaning use whatever is defined in the image.
	// +optional
	Command []string `json:"command,omitempty" patchStrategy:"replace"`

	// The arguments to supply to the command running the dockerimage component. The arguments are supplied either to the default command provided in the image or to the overridden command.
	//
	// Defaults to an empty array, meaning use whatever is defined in the image.
	// +optional
	Args []string `json:"args,omitempty" patchStrategy:"replace"`

	// Toggles whether or not the project source code should
	// be mounted in the component.
	//
	// Defaults to true for all component types except plugins and components that set `dedicatedPod` to true.
	// +optional
	MountSources *bool `json:"mountSources,omitempty"`

	// Optional specification of the path in the container where
	// project sources should be transferred/mounted when `mountSources` is `true`.
	// When omitted, the default value of /projects is used.
	// +optional
	SourceMapping string `json:"sourceMapping,omitempty"`

	// Specify if a container should run in its own separated pod,
	// instead of running as part of the main development environment pod.
	//
	// Default value is `false`
	// +optional
	DedicatedPod bool `json:"dedicatedPod,omitempty"`
}

type EndpointPluginOverride struct {

	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`

	//  +optional
	TargetPort int `json:"targetPort,omitempty"`

	// Describes how the endpoint should be exposed on the network.
	//
	// - `public` means that the endpoint will be exposed on the public network, typically through
	// a K8S ingress or an OpenShift route.
	//
	// - `internal` means that the endpoint will be exposed internally outside of the main workspace POD,
	// typically by K8S services, to be consumed by other elements running
	// on the same cloud internal network.
	//
	// - `none` means that the endpoint will not be exposed and will only be accessible
	// inside the main workspace POD, on a local address.
	//
	// Default value is `public`
	// +optional
	Exposure EndpointExposurePluginOverride `json:"exposure,omitempty"`

	// Describes the application and transport protocols of the traffic that will go through this endpoint.
	//
	// - `http`: Endpoint will have `http` traffic, typically on a TCP connection.
	// It will be automaticaly promoted to `https` when the `secure` field is set to `true`.
	//
	// - `https`: Endpoint will have `https` traffic, typically on a TCP connection.
	//
	// - `ws`: Endpoint will have `ws` traffic, typically on a TCP connection.
	// It will be automaticaly promoted to `wss` when the `secure` field is set to `true`.
	//
	// - `wss`: Endpoint will have `wss` traffic, typically on a TCP connection.
	//
	// - `tcp`: Endpoint will have traffic on a TCP connection, without specifying an application protocol.
	//
	// - `udp`: Endpoint will have traffic on an UDP connection, without specifying an application protocol.
	//
	// Default value is `http`
	// +optional
	Protocol EndpointProtocolPluginOverride `json:"protocol,omitempty"`

	// Describes whether the endpoint should be secured and protected by some
	// authentication process. This requires a protocol of `https` or `wss`.
	// +optional
	Secure bool `json:"secure,omitempty"`

	// Path of the endpoint URL
	// +optional
	Path string `json:"path,omitempty"`

	// Map of implementation-dependant string-based free-form attributes.
	//
	// Examples of Che-specific attributes:
	//
	// - cookiesAuthEnabled: "true" / "false",
	//
	// - type: "terminal" / "ide" / "ide-dev",
	// +optional
	Attributes attributes.Attributes `json:"attributes,omitempty"`
}

type K8sLikeComponentPluginOverride struct {
	BaseComponentPluginOverride            `json:",inline"`
	K8sLikeComponentLocationPluginOverride `json:",inline"`
	Endpoints                              []EndpointPluginOverride `json:"endpoints,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// Volume that should be mounted to a component container
type VolumePluginOverride struct {

	// +optional
	// Size of the volume
	Size string `json:"size,omitempty"`

	// +optional
	// Ephemeral volumes are not stored persistently across restarts. Defaults
	// to false
	Ephemeral bool `json:"ephemeral,omitempty"`
}

type LabeledCommandPluginOverride struct {
	BaseCommandPluginOverride `json:",inline"`

	// +optional
	// Optional label that provides a label for this command
	// to be used in Editor UI menus for example
	Label string `json:"label,omitempty"`
}

type EnvVarPluginOverride struct {
	Name string `json:"name" yaml:"name"`
	//  +optional
	Value string `json:"value,omitempty" yaml:"value"`
}

// Volume that should be mounted to a component container
type VolumeMountPluginOverride struct {

	// The volume mount name is the name of an existing `Volume` component.
	// If several containers mount the same volume name
	// then they will reuse the same volume and will be able to access to the same files.
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`

	// The path in the component container where the volume should be mounted.
	// If not path is mentioned, default path is the is `/<name>`.
	// +optional
	Path string `json:"path,omitempty"`
}

// EndpointExposure describes the way an endpoint is exposed on the network.
// Only one of the following exposures may be specified: public, internal, none.
// +kubebuilder:validation:Enum=public;internal;none
type EndpointExposurePluginOverride string

// EndpointProtocol defines the application and transport protocols of the traffic that will go through this endpoint.
// Only one of the following protocols may be specified: http, ws, tcp, udp.
// +kubebuilder:validation:Enum=http;https;ws;wss;tcp;udp
type EndpointProtocolPluginOverride string

// +union
type K8sLikeComponentLocationPluginOverride struct {

	// +kubebuilder:validation:Enum=Uri;Inlined
	// Type of Kubernetes-like location
	// +
	// +unionDiscriminator
	// +optional
	LocationType K8sLikeComponentLocationTypePluginOverride `json:"locationType,omitempty"`

	// Location in a file fetched from a uri.
	// +optional
	Uri string `json:"uri,omitempty"`

	// Inlined manifest
	// +optional
	Inlined string `json:"inlined,omitempty"`
}

type BaseCommandPluginOverride struct {

	// +optional
	// Defines the group this command is part of
	Group *CommandGroupPluginOverride `json:"group,omitempty"`
}

// K8sLikeComponentLocationType describes the type of
// the location the configuration is fetched from.
// Only one of the following component type may be specified.
type K8sLikeComponentLocationTypePluginOverride string

type CommandGroupPluginOverride struct {

	//  +optional
	// Kind of group the command is part of
	Kind CommandGroupKindPluginOverride `json:"kind,omitempty"`

	// +optional
	// Identifies the default command for a given group kind
	IsDefault bool `json:"isDefault,omitempty"`
}

// CommandGroupKind describes the kind of command group.
// +kubebuilder:validation:Enum=build;run;test;debug
type CommandGroupKindPluginOverride string

func (overrides PluginOverrides) isOverride() {}
