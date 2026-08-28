import { httpClient } from "../core";
import type { CreateRolePayload, Role } from "./types";

export async function listRoles(): Promise<Role[]> {
  const { data } = await httpClient.get<Role[]>("/roles");
  return data;
}

export async function createRole(payload: CreateRolePayload): Promise<Role> {
  const { data } = await httpClient.post<Role>("/roles", payload);
  return data;
}

export async function updateRole(id: string, payload: Partial<CreateRolePayload>): Promise<Role> {
  const { data } = await httpClient.patch<Role>(`/roles/${id}`, payload);
  return data;
}

export async function listPermissions(): Promise<string[]> {
  const { data } = await httpClient.get<string[]>("/permissions");
  return data;
}

export type * from "./types";
