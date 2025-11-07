// Auto-generated from schema/nebulabox.schema.json
// DO NOT EDIT MANUALLY
package generated

import "time"

type Container struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Status string `json:"status"`
	Created time.Time `json:"created"`
	Started *time.Time `json:"started,omitempty"`
	Stopped *time.Time `json:"stopped,omitempty"`
	Command *string `json:"command,omitempty"`
	Env []string `json:"env,omitempty"`
	Ports []string `json:"ports,omitempty"`
	Volumes []string `json:"volumes,omitempty"`
	Network *string `json:"network,omitempty"`
	WorkspaceId *string `json:"workspaceId,omitempty"`
	Labels interface{} `json:"labels,omitempty"`
}

type Image struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Tag string `json:"tag"`
	Digest *string `json:"digest,omitempty"`
	Size string `json:"size"`
	Created time.Time `json:"created"`
	Registry *string `json:"registry,omitempty"`
	Repository *string `json:"repository,omitempty"`
}

type Workspace struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
	Status string `json:"status"`
	OwnerId string `json:"ownerId"`
	ContainerId *string `json:"containerId,omitempty"`
	Members []string `json:"members,omitempty"`
	Settings interface{} `json:"settings,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WorkspaceMember struct {
	UserId string `json:"userId"`
	Username string `json:"username,omitempty"`
	Role string `json:"role"`
	JoinedAt time.Time `json:"joinedAt"`
	LastSeen *time.Time `json:"lastSeen,omitempty"`
	IsActive bool `json:"isActive,omitempty"`
}

type WorkspaceSettings struct {
	AllowGuestAccess bool `json:"allowGuestAccess,omitempty"`
	MaxMembers int `json:"maxMembers,omitempty"`
	SessionTimeout int `json:"sessionTimeout,omitempty"`
	AutoPauseOnIdle bool `json:"autoPauseOnIdle,omitempty"`
	IdleTimeout int `json:"idleTimeout,omitempty"`
	AuditLogging bool `json:"auditLogging,omitempty"`
}
