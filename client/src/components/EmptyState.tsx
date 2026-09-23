import type { ReactNode } from "react";
import "./EmptyState.css";

interface EmptyStateProps {
  title: string;
  subtitle?: string;
  action?: ReactNode;
}

export function EmptyState({ title, subtitle, action }: EmptyStateProps) {
  return (
    <div className="empty-state">
      <p className="empty-state-title">{title}</p>
      {subtitle && <p className="empty-state-subtitle">{subtitle}</p>}
      {action}
    </div>
  );
}
