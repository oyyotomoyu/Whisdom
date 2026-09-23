import { request } from "../core";
import type { CreateUserPayload, UpdateUserPayload, UserAccount } from "./types";

export async function listUsers(): Promise<UserAccount[]> {
  return request<UserAccount[]>("/users");
}

export async function getUser(id: string): Promise<UserAccount> {
  return request<UserAccount>(`/users/${id}`);
}

export async function createUser(payload: CreateUserPayload): Promise<UserAccount> {
  return request<UserAccount>("/users", { method: "POST", body: JSON.stringify(payload) });
}

export async function updateUser(id: string, payload: UpdateUserPayload): Promise<UserAccount> {
  return request<UserAccount>(`/users/${id}`, { method: "PATCH", body: JSON.stringify(payload) });
}

export async function deleteUser(id: string): Promise<void> {
  await request<void>(`/users/${id}`, { method: "DELETE" });
}

export type * from "./types";
