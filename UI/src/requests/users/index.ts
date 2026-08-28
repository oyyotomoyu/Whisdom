import { httpClient } from "../core";
import type { CreateUserPayload, UpdateUserPayload, UserAccount } from "./types";

export async function listUsers(): Promise<UserAccount[]> {
  const { data } = await httpClient.get<UserAccount[]>("/users");
  return data;
}

export async function getUser(id: string): Promise<UserAccount> {
  const { data } = await httpClient.get<UserAccount>(`/users/${id}`);
  return data;
}

export async function createUser(payload: CreateUserPayload): Promise<UserAccount> {
  const { data } = await httpClient.post<UserAccount>("/users", payload);
  return data;
}

export async function updateUser(id: string, payload: UpdateUserPayload): Promise<UserAccount> {
  const { data } = await httpClient.patch<UserAccount>(`/users/${id}`, payload);
  return data;
}

export async function deleteUser(id: string): Promise<void> {
  await httpClient.delete(`/users/${id}`);
}

export type * from "./types";
