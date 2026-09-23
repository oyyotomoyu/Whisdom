import { request } from "../core";
import type { ModelInfo } from "./types";

export async function listModels(): Promise<ModelInfo[]> {
  return request<ModelInfo[]>("/models");
}

export async function getCurrentModel(): Promise<ModelInfo> {
  return request<ModelInfo>("/models/current");
}

export async function activateModel(id: string): Promise<void> {
  await request<void>(`/models/${id}/activate`, { method: "POST" });
}

export type * from "./types";
