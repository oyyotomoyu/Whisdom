import { clearAccessToken, request, setAccessToken } from "../core";
import type { AuthUser } from "../../store/types";
import { getDevCurrentUser, getDevLoginResponse, isDevSession } from "./devMode";

export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  user: AuthUser;
}

export async function login(payload: LoginPayload): Promise<LoginResponse> {
  const devResponse = getDevLoginResponse(payload);
  if (devResponse) {
    setAccessToken(devResponse.access_token);
    return devResponse;
  }

  const response = await request<LoginResponse>(
    "/auth/login",
    { method: "POST", body: JSON.stringify(payload) },
    { auth: false, retryOnUnauthorized: false }
  );
  setAccessToken(response.access_token);
  return response;
}

export async function logout(): Promise<void> {
  if (isDevSession()) {
    clearAccessToken();
    return;
  }

  try {
    await request<void>("/auth/logout", { method: "POST" });
  } finally {
    clearAccessToken();
  }
}

export async function fetchCurrentUser(): Promise<AuthUser> {
  const devUser = getDevCurrentUser();
  if (devUser) {
    return devUser;
  }

  return request<AuthUser>("/auth/me");
}
