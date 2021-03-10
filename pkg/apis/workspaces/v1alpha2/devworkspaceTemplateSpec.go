package v1alpha2

// Structure of the workspace. This is also the specification of a workspace template.
// +devfile:jsonschema:generate
type DevWorkspaceTemplateSpec struct {
	// Parent workspace template
	// +optional
	Parent *Parent `json:"parent,omitempty"`

	DevWorkspaceTemplateSpecContent `json:",inline"`
}

// +devfile:overrides:generate
type DevWorkspaceTemplateSpecContent struct {
	// Map of implementation-dependant free-form YAML attributes.
	// Attribute values can be referenced throughout the devfile in string type fields in the form {{attribute-key}}
	// except for schemaVersion, metadata, parent source. Exception to the string field also include element's key identifiers
	// (command id, component name, endpoint name, project name, etc.) and their references(events, command's component, container's
	// volume mount name, etc.) and string enums(command group kind, endpoint exposure, etc.)
	// +optional
	// +patchStrategy=merge
	// +devfile:overrides:include:omitInPlugin=true,description=Overrides of attributes encapsulated in a parent devfile.
	Attributes map[string]string `json:"attributes,omitempty" patchStrategy:"merge"`

	// List of the workspace components, such as editor and plugins,
	// user-provided containers, or other types of components
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +devfile:overrides:include:description=Overrides of components encapsulated in a parent devfile or a plugin.
	// +devfile:toplevellist
	Components []Component `json:"components,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Projects worked on in the workspace, containing names and sources locations
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +devfile:overrides:include:omitInPlugin=true,description=Overrides of projects encapsulated in a parent devfile.
	// +devfile:toplevellist
	Projects []Project `json:"projects,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// StarterProjects is a project that can be used as a starting point when bootstrapping new projects
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +devfile:overrides:include:omitInPlugin=true,description=Overrides of starterProjects encapsulated in a parent devfile.
	// +devfile:toplevellist
	StarterProjects []StarterProject `json:"starterProjects,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Predefined, ready-to-use, workspace-related commands
	// +optional
	// +patchMergeKey=id
	// +patchStrategy=merge
	// +devfile:overrides:include:description=Overrides of commands encapsulated in a parent devfile or a plugin.
	// +devfile:toplevellist
	Commands []Command `json:"commands,omitempty" patchStrategy:"merge" patchMergeKey:"id"`

	// Bindings of commands to events.
	// Each command is referred-to by its name.
	// +optional
	// +devfile:overrides:include:omit=true
	Events *Events `json:"events,omitempty"`
}
