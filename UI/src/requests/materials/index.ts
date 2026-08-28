import { httpClient } from "../core";
import type { Material, MaterialProcessingStatus } from "./types";

export async function listMaterials(): Promise<Material[]> {
  const { data } = await httpClient.get<Material[]>("/materials");
  return data;
}

export async function getMaterial(id: string): Promise<Material> {
  const { data } = await httpClient.get<Material>(`/materials/${id}`);
  return data;
}

export interface UploadMaterialPayload {
  file: File;
  destinationPath: string;
  onProgress?: (percent: number) => void;
}

export async function uploadMaterial({ file, destinationPath, onProgress }: UploadMaterialPayload): Promise<Material> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("destination_path", destinationPath);

  const { data } = await httpClient.post<Material>("/materials", formData, {
    headers: { "Content-Type": "multipart/form-data" },
    onUploadProgress: (event) => {
      if (onProgress && event.total) {
        onProgress(Math.round((event.loaded / event.total) * 100));
      }
    },
  });
  return data;
}

export async function updateMaterialDestination(id: string, destinationPath: string): Promise<Material> {
  const { data } = await httpClient.patch<Material>(`/materials/${id}`, { destination_path: destinationPath });
  return data;
}

export async function deleteMaterial(id: string): Promise<void> {
  await httpClient.delete(`/materials/${id}`);
}

export async function reprocessMaterial(id: string): Promise<void> {
  await httpClient.post(`/materials/${id}/process`);
}

export async function getMaterialStatus(id: string): Promise<MaterialProcessingStatus> {
  const { data } = await httpClient.get<MaterialProcessingStatus>(`/materials/${id}/status`);
  return data;
}

export function validateDestinationPath(path: string): string | null {
  if (!path.trim()) return "required";
  if (!path.startsWith("/")) return "leadingSlash";
  if (path.includes("..")) return "parentTraversal";
  return null;
}

export type * from "./types";
