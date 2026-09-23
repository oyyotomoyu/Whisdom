export const API_BASE_URL = "/api/v1";

export interface ApiErrorBody {
  error?: string;
  message?: string;
  code?: string;
  [key: string]: unknown;
}

export class ApiError extends Error {
  readonly status: number;
  readonly body: ApiErrorBody | null;
  readonly code?: string;

  constructor(status: number, body: ApiErrorBody | null) {
    super(body?.error ?? body?.message ?? `Request failed with status ${status}`);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
    this.code = body?.code;
  }

  get isUnauthorized() {
    return this.status === 401;
  }

  get isForbidden() {
    return this.status === 403;
  }

  get isNetworkError() {
    return this.status === 0;
  }
}

export interface RequestOptions {
  auth?: boolean;
  retryOnUnauthorized?: boolean;
}

let accessToken: string | null = null;
let refreshPromise: Promise<string | null> | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getAccessToken() {
  return accessToken;
}

export function clearAccessToken() {
  accessToken = null;
}

export async function request<TResponse>(
  path: string,
  init: RequestInit = {},
  options: RequestOptions = {}
): Promise<TResponse> {
  const { auth = true, retryOnUnauthorized = true } = options;

  let response: Response;
  try {
    response = await send(path, init, auth);
  } catch {
    throw new ApiError(0, { error: "Network error" });
  }

  if (response.status === 401 && auth && retryOnUnauthorized) {
    const refreshedToken = await refreshAccessToken();
    if (refreshedToken) {
      return request<TResponse>(path, init, { auth, retryOnUnauthorized: false });
    }
  }

  return parseResponse<TResponse>(response);
}

async function send(path: string, init: RequestInit, auth: boolean) {
  const headers = new Headers(init.headers);
  const isFormData = typeof FormData !== "undefined" && init.body instanceof FormData;

  if (init.body != null && !isFormData && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (auth && accessToken) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  }

  return fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers,
    credentials: "include",
  });
}

async function refreshAccessToken(): Promise<string | null> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const response = await send("/auth/refresh", { method: "POST" }, false);
        if (!response.ok) {
          clearAccessToken();
          return null;
        }
        const tokens = (await response.json()) as { access_token: string };
        setAccessToken(tokens.access_token);
        return tokens.access_token;
      } catch {
        clearAccessToken();
        return null;
      }
    })().finally(() => {
      refreshPromise = null;
    });
  }

  return refreshPromise;
}

async function parseResponse<TResponse>(response: Response): Promise<TResponse> {
  if (!response.ok) {
    throw new ApiError(response.status, await readErrorBody(response));
  }
  if (response.status === 204) {
    return undefined as TResponse;
  }

  const text = await response.text();
  if (!text) {
    return undefined as TResponse;
  }
  return JSON.parse(text) as TResponse;
}

async function readErrorBody(response: Response): Promise<ApiErrorBody | null> {
  try {
    return (await response.json()) as ApiErrorBody;
  } catch {
    return null;
  }
}
