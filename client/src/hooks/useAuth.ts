import { useCallback } from "react";
import { useAuthStore } from "../store/authStore";
import { login as loginRequest, logout as logoutRequest } from "../requests/auth";
import { ApiError } from "../requests/core";

export function useAuth() {
  const user = useAuthStore((state) => state.user);
  const accessToken = useAuthStore((state) => state.accessToken);
  const isAuthLoading = useAuthStore((state) => state.isAuthLoading);
  const hasPermission = useAuthStore((state) => state.hasPermission);
  const isAdmin = useAuthStore((state) => state.isAdmin);
  const setSession = useAuthStore((state) => state.setSession);
  const clearSession = useAuthStore((state) => state.clearSession);
  const setAuthLoading = useAuthStore((state) => state.setAuthLoading);

  const login = useCallback(
    async (username: string, password: string) => {
      try {
        const { access_token, user } = await loginRequest({ username, password });
        setSession(user, access_token);
        return { ok: true as const };
      } catch (error) {
        const message = error instanceof ApiError ? error.message : "login.invalidCredentials";
        return { ok: false as const, message };
      }
    },
    [setSession]
  );

  const logout = useCallback(async () => {
    try {
      await logoutRequest();
    } finally {
      clearSession();
    }
  }, [clearSession]);

  return {
    user,
    accessToken,
    isAuthenticated: Boolean(user && accessToken),
    isAuthLoading,
    hasPermission,
    isAdmin: isAdmin(),
    login,
    logout,
    setAuthLoading,
  };
}
