import { httpClient } from "../core";
import type { ModelInfo } from "./types";

export async function listModels(): Promise<ModelInfo[]> {
  const { data } = await httpClient.get<ModelInfo[]>("/models");
  return data;
}

export async function getCurrentModel(): Promise<ModelInfo> {
  const { data } = await httpClient.get<ModelInfo>("/models/current");
  return data;
}

export async function activateModel(id: string): Promise<void> {
  await httpClient.post(`/models/${id}/activate`);
}

export type * from "./types";
