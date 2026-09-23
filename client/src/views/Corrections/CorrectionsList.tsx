import { useTranslation } from "../../../node_modules/react-i18next";
import type { Correction } from "../../requests/corrections/types";
import "./CorrectionsList.css";

interface CorrectionsListProps {
  corrections: Correction[];
  onDelete: (id: string) => void;
}

export function CorrectionsList({ corrections, onDelete }: CorrectionsListProps) {
  const { t } = useTranslation("corrections");

  if (corrections.length === 0) {
    return <p className="corrections-empty">{t("list.empty")}</p>;
  }

  return (
    <ul className="corrections-list">
      {corrections.map((correction) => (
        <li key={correction.id} className="correction-card">
          <div className="correction-card-header">
            <span className={`status-pill status-${correction.status}`}>{t(`status.${correction.status}`)}</span>
            <span className="correction-date">{new Date(correction.createdAt).toLocaleString()}</span>
            <button type="button" className="correction-delete" onClick={() => onDelete(correction.id)}>
              ×
            </button>
          </div>
          <p className="correction-question">{correction.question}</p>
          <p className="correction-original">{correction.originalAnswer}</p>
          <p className="correction-corrected">{correction.correctedAnswer}</p>
          {correction.note && <p className="correction-note">{correction.note}</p>}
        </li>
      ))}
    </ul>
  );
}
