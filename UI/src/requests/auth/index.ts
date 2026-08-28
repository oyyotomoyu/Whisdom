import { httpClient } from "../core";
import type { AuthUser } from "../../store/types";

export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  user: AuthUser;
}

export async function login(payload: LoginPayload): Promise<LoginResponse> {
  const { data } = await httpClient.post<LoginResponse>("/auth/login", payload);
  return data;
}

export async function logout(): Promise<void> {
  await httpClient.post("/auth/logout");
}

export async function fetchCurrentUser(): Promise<AuthUser> {
  const { data } = await httpClient.get<AuthUser>("/auth/me");
  return data;
}
