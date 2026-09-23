import type { AuthUser, Permission } from "../../store/types";
import type { LoginPayload, LoginResponse } from "./index";

const devAccessToken = "dev-access-token";
const devEmail = "admin@whisdom.local";
const devPassword = "1234";

const devUser: AuthUser = {
  id: "dev-admin",
  name: "Development Admin",
  email: devEmail,
  role: "administrator",
  permissions: [
    "chat.use",
    "materials.read",
    "materials.upload",
    "materials.delete",
    "corrections.create",
    "corrections.read",
    "corrections.delete",
    "users.read",
    "users.manage",
    "models.read",
    "models.manage",
    "system.manage",
  ] satisfies Permission[],
};

export function isDevMode() {
  if (process.env.NODE_ENV !== "production") {
    return true;
  }

  return false;
}

export function getDevCredentials() {
  return {
    username: devEmail,
    password: devPassword,
  };
}

export function getDevLoginResponse(payload: LoginPayload): LoginResponse | null {
  if (process.env.NODE_ENV !== "production") {
    if (payload.username === devEmail && payload.password === devPassword) {
      return {
        access_token: devAccessToken,
        user: devUser,
      };
    }
  }

  return null;
}

export function getDevCurrentUser(): AuthUser | null {
  if (process.env.NODE_ENV !== "production") {
    if (localStorage.getItem("whisdom_auth")?.includes(devAccessToken)) {
      return devUser;
    }
  }

  return null;
}

export function isDevSession() {
  return getDevCurrentUser() != null;
}
