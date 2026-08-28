import { httpClient } from "../core";
import type { Correction, CreateCorrectionPayload } from "./types";

export async function listCorrections(): Promise<Correction[]> {
  const { data } = await httpClient.get<Correction[]>("/corrections");
  return data;
}

export async function getCorrection(id: string): Promise<Correction> {
  const { data } = await httpClient.get<Correction>(`/corrections/${id}`);
  return data;
}

export async function createCorrection(payload: CreateCorrectionPayload): Promise<Correction> {
  const endpoint = payload.messageId ? `/messages/${payload.messageId}/corrections` : "/corrections";
  const { data } = await httpClient.post<Correction>(endpoint, payload);
  return data;
}

export async function deleteCorrection(id: string): Promise<void> {
  await httpClient.delete(`/corrections/${id}`);
}

export type * from "./types";
