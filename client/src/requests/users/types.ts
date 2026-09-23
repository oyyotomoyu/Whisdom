export interface UserAccount {
  id: string;
  name: string;
  email: string;
  role: string;
  active: boolean;
}

export interface CreateUserPayload {
  name: string;
  email: string;
  role: string;
  password: string;
}

export type UpdateUserPayload = Partial<Pick<CreateUserPayload, "name" | "email" | "role">> & { active?: boolean };
