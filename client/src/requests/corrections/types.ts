export type CorrectionUsage = "rag" | "evaluation" | "training" | "all";
export type CorrectionStatus = "pending" | "approved" | "rejected";

export interface Correction {
  id: string;
  question: string;
  originalAnswer: string;
  correctedAnswer: string;
  note?: string;
  relatedMaterialId?: string;
  usage: CorrectionUsage;
  status: CorrectionStatus;
  createdAt: string;
}

export interface CreateCorrectionPayload {
  question: string;
  originalAnswer: string;
  correctedAnswer: string;
  note?: string;
  relatedMaterialId?: string;
  usage: CorrectionUsage;
  messageId?: string;
}
