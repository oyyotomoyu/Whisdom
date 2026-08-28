export interface Role {
  id: string;
  name: string;
  permissions: string[];
}

export interface CreateRolePayload {
  name: string;
  permissions: string[];
}
