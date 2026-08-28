import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "../../store/authStore";
import { ApiError } from "./apiError";

export const API_BASE_URL = "/api/v1";

export const httpClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
});

httpClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken;
  if (token) {
    config.headers.set("Authorization", `Bearer ${token}`);
  }
  return config;
});

let refreshPromise: Promise<string | null> | null = null;

async function refreshAccessToken(): Promise<string | null> {
  if (!refreshPromise) {
    refreshPromise = axios
      .post<{ access_token: string }>(`${API_BASE_URL}/auth/refresh`, {}, { withCredentials: true })
      .then((res) => res.data.access_token)
      .catch(() => null)
      .finally(() => {
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

httpClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<{ message?: string; code?: string }>) => {
    const originalRequest = error.config as (InternalAxiosRequestConfig & { _retried?: boolean }) | undefined;

    if (!error.response) {
      return Promise.reject(new ApiError("Network error", 0));
    }

    const { status, data } = error.response;

    if (status === 401 && originalRequest && !originalRequest._retried && !originalRequest.url?.includes("/auth/")) {
      originalRequest._retried = true;
      const newToken = await refreshAccessToken();
      if (newToken) {
        const user = useAuthStore.getState().user;
        if (user) {
          useAuthStore.getState().setSession(user, newToken);
        }
        return httpClient(originalRequest);
      }
      useAuthStore.getState().clearSession();
    }

    return Promise.reject(new ApiError(data?.message ?? "Request failed", status, data?.code));
  }
);
