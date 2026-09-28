import type { AuthUser, Permission } from "../../store/types";
import type { LoginPayload, LoginResponse } from "./index";

interface DevAccount {
  accessToken: string;
  email: string;
  password: string;
  user: AuthUser;
}

const devAccounts: DevAccount[] = [
  {
    accessToken: "dev-access-token-admin",
    email: "admin@whisdom.local",
    password: "1234",
    user: {
      id: "dev-admin",
      name: "Development Admin",
      email: "admin@whisdom.local",
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
    },
  },
  {
    accessToken: "dev-access-token-user",
    email: "user@whisdom.local",
    password: "1234",
    user: {
      id: "dev-user",
      name: "Development User",
      email: "user@whisdom.local",
      role: "user",
      permissions: ["chat.use"] satisfies Permission[],
    },
  },
];

export function isDevMode() {
  if (process.env.NODE_ENV !== "production") {
    return true;
  }

  return false;
}

export function getDevCredentials() {
  const [admin] = devAccounts;
  return {
    username: admin.email,
    password: admin.password,
  };
}

export function getDevLoginResponse(payload: LoginPayload): LoginResponse | null {
  if (process.env.NODE_ENV !== "production") {
    const account = devAccounts.find(
      (candidate) => candidate.email === payload.username && candidate.password === payload.password
    );

    if (account) {
      return {
        access_token: account.accessToken,
        user: account.user,
      };
    }
  }

  return null;
}

export function getDevCurrentUser(): AuthUser | null {
  if (process.env.NODE_ENV !== "production") {
    const stored = localStorage.getItem("whisdom_auth");
    const account = stored ? devAccounts.find((candidate) => stored.includes(candidate.accessToken)) : undefined;

    if (account) {
      return account.user;
    }
  }

  return null;
}

export function isDevSession() {
  return getDevCurrentUser() != null;
}
