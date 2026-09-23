import { request } from "../core";
import type { Correction, CreateCorrectionPayload } from "./types";

export async function listCorrections(): Promise<Correction[]> {
  return request<Correction[]>("/corrections");
}

export async function getCorrection(id: string): Promise<Correction> {
  return request<Correction>(`/corrections/${id}`);
}

export async function createCorrection(payload: CreateCorrectionPayload): Promise<Correction> {
  const endpoint = payload.messageId ? `/messages/${payload.messageId}/corrections` : "/corrections";
  return request<Correction>(endpoint, { method: "POST", body: JSON.stringify(payload) });
}

export async function deleteCorrection(id: string): Promise<void> {
  await request<void>(`/corrections/${id}`, { method: "DELETE" });
}

export type * from "./types";
