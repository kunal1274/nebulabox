// Auto-generated from schema/nebulabox.schema.json
// DO NOT EDIT MANUALLY

export interface Container {
  id: string;
  name: string;
  image: string;
  status: string;
  created: string;
  started?: string | null;
  stopped?: string | null;
  command?: string | null;
  env?: Array<string>;
  ports?: Array<string>;
  volumes?: Array<string>;
  network?: string | null;
  workspaceId?: string | null;
  labels?: Record<string, any>;
}

export interface Image {
  id: string;
  name: string;
  tag: string;
  digest?: string | null;
  size: string;
  created: string;
  registry?: string | null;
  repository?: string | null;
}

export interface Workspace {
  id: string;
  name: string;
  description?: string | null;
  status: string;
  ownerId: string;
  containerId?: string | null;
  members?: Array<any>;
  settings?: any;
  metadata?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface WorkspaceMember {
  userId: string;
  username?: string;
  role: string;
  joinedAt: string;
  lastSeen?: string | null;
  isActive?: boolean;
}

export interface WorkspaceSettings {
  allowGuestAccess?: boolean;
  maxMembers?: number;
  sessionTimeout?: number;
  autoPauseOnIdle?: boolean;
  idleTimeout?: number;
  auditLogging?: boolean;
}
