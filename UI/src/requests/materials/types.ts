export type MaterialStatus = "pending" | "processing" | "ready" | "failed";

export interface Material {
  id: string;
  filename: string;
  type: string;
  uploadedAt: string;
  uploadedBy: string;
  status: MaterialStatus;
  destinationPath: string;
  ragAvailable: boolean;
  trainingAvailable: boolean;
}

export interface MaterialProcessingStatus {
  id: string;
  filename: string;
  status: MaterialStatus;
  processing: {
    extraction: string;
    chunking: string;
    embedding: string;
  };
}
