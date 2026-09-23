export type Permission =
  | "chat.use"
  | "materials.read"
  | "materials.upload"
  | "materials.delete"
  | "corrections.create"
  | "corrections.read"
  | "corrections.delete"
  | "users.read"
  | "users.manage"
  | "models.read"
  | "models.manage"
  | "system.manage";

export interface AuthUser {
  id: string;
  name: string;
  email: string;
  role: string;
  permissions: Permission[];
}
