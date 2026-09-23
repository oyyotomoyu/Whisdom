import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { AuthUser, Permission } from "./types";

interface AuthState {
  user: AuthUser | null;
  accessToken: string | null;
  isAuthLoading: boolean;
  setSession: (user: AuthUser, accessToken: string) => void;
  clearSession: () => void;
  setAuthLoading: (loading: boolean) => void;
  hasPermission: (permission: Permission) => boolean;
  isAdmin: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      isAuthLoading: true,
      setSession: (user, accessToken) => set({ user, accessToken, isAuthLoading: false }),
      clearSession: () => set({ user: null, accessToken: null, isAuthLoading: false }),
      setAuthLoading: (loading) => set({ isAuthLoading: loading }),
      hasPermission: (permission) => get().user?.permissions.includes(permission) ?? false,
      isAdmin: () => get().user?.role === "administrator",
    }),
    {
      name: "whisdom_auth",
      partialize: (state) => ({ user: state.user, accessToken: state.accessToken }),
    }
  )
);
