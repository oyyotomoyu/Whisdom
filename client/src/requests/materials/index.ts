import { request } from "../core";
import type { Material, MaterialProcessingStatus, MaterialSourceType } from "./types";

export async function listMaterials(): Promise<Material[]> {
  return request<Material[]>("/materials");
}

export async function getMaterial(id: string): Promise<Material> {
  return request<Material>(`/materials/${id}`);
}

export interface UploadMaterialPayload {
  file: File;
  destinationPath: string;
  sourceType?: MaterialSourceType;
  sourceLocation?: string;
  onProgress?: (percent: number) => void;
}

export interface AddMaterialSourcePayload {
  sourceType: Exclude<MaterialSourceType, "upload">;
  sourceLocation: string;
  destinationPath: string;
}

export async function uploadMaterial({
  file,
  destinationPath,
  sourceType = "upload",
  sourceLocation,
  onProgress,
}: UploadMaterialPayload): Promise<Material> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("destination_path", destinationPath);
  formData.append("source_type", sourceType);
  if (sourceLocation) formData.append("source_location", sourceLocation);

  onProgress?.(0);
  const material = await request<Material>("/materials", { method: "POST", body: formData });
  onProgress?.(100);
  return material;
}

export async function addMaterialSource(payload: AddMaterialSourcePayload): Promise<Material> {
  return request<Material>("/materials", {
    method: "POST",
    body: JSON.stringify({
      source_type: payload.sourceType,
      source_location: payload.sourceLocation,
      destination_path: payload.destinationPath,
    }),
  });
}

export async function updateMaterialDestination(id: string, destinationPath: string): Promise<Material> {
  return request<Material>(`/materials/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ destination_path: destinationPath }),
  });
}

export async function deleteMaterial(id: string): Promise<void> {
  await request<void>(`/materials/${id}`, { method: "DELETE" });
}

export async function reprocessMaterial(id: string): Promise<void> {
  await request<void>(`/materials/${id}/process`, { method: "POST" });
}

export async function getMaterialStatus(id: string): Promise<MaterialProcessingStatus> {
  return request<MaterialProcessingStatus>(`/materials/${id}/status`);
}

export function validateDestinationPath(path: string): string | null {
  if (!path.trim()) return "required";
  if (!path.startsWith("/")) return "leadingSlash";
  if (path.includes("..")) return "parentTraversal";
  return null;
}

export function validateSourceLocation(sourceType: MaterialSourceType, sourceLocation: string): string | null {
  const value = sourceLocation.trim();
  if (sourceType === "upload") return null;
  if (!value) return "required";
  if ((sourceType === "local_path" || value.startsWith("/")) && value.includes("..")) return "parentTraversal";
  if (sourceType !== "local_path") {
    try {
      const url = new URL(value);
      if (url.protocol !== "https:" && url.protocol !== "http:") return "scheme";
    } catch {
      return "url";
    }
  }
  return null;
}

export type * from "./types";
