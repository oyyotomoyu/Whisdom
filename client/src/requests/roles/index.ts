import { request } from "../core";
import type { CreateRolePayload, Role } from "./types";

export async function listRoles(): Promise<Role[]> {
  return request<Role[]>("/roles");
}

export async function createRole(payload: CreateRolePayload): Promise<Role> {
  return request<Role>("/roles", { method: "POST", body: JSON.stringify(payload) });
}

export async function updateRole(id: string, payload: Partial<CreateRolePayload>): Promise<Role> {
  return request<Role>(`/roles/${id}`, { method: "PATCH", body: JSON.stringify(payload) });
}

export async function listPermissions(): Promise<string[]> {
  return request<string[]>("/permissions");
}

export type * from "./types";
