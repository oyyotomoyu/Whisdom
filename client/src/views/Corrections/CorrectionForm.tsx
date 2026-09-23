import { useState, type FormEvent } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import type { CorrectionUsage } from "../../requests/corrections/types";
import "./CorrectionForm.css";

interface CorrectionFormProps {
  initialAnswer?: string;
  initialMessageId?: string;
  onSubmit: (values: {
    question: string;
    originalAnswer: string;
    correctedAnswer: string;
    note: string;
    relatedMaterialId: string;
    usage: CorrectionUsage;
    messageId?: string;
  }) => Promise<void>;
}

const usageOptions: CorrectionUsage[] = ["rag", "evaluation", "training", "all"];

export function CorrectionForm({ initialAnswer, initialMessageId, onSubmit }: CorrectionFormProps) {
  const { t } = useTranslation("corrections");

  const [question, setQuestion] = useState("");
  const [originalAnswer, setOriginalAnswer] = useState(initialAnswer ?? "");
  const [correctedAnswer, setCorrectedAnswer] = useState("");
  const [note, setNote] = useState("");
  const [relatedMaterialId, setRelatedMaterialId] = useState("");
  const [usage, setUsage] = useState<CorrectionUsage>("all");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setIsSubmitting(true);
    try {
      await onSubmit({
        question,
        originalAnswer,
        correctedAnswer,
        note,
        relatedMaterialId,
        usage,
        messageId: initialMessageId,
      });
      setQuestion("");
      setOriginalAnswer("");
      setCorrectedAnswer("");
      setNote("");
      setRelatedMaterialId("");
      setUsage("all");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form className="correction-form" onSubmit={handleSubmit}>
      <label>{t("form.question")}</label>
      <textarea value={question} onChange={(e) => setQuestion(e.target.value)} required rows={2} />

      <label>{t("form.originalAnswer")}</label>
      <textarea value={originalAnswer} onChange={(e) => setOriginalAnswer(e.target.value)} required rows={2} />

      <label>{t("form.correctedAnswer")}</label>
      <textarea value={correctedAnswer} onChange={(e) => setCorrectedAnswer(e.target.value)} required rows={3} />

      <label>{t("form.note")}</label>
      <textarea value={note} onChange={(e) => setNote(e.target.value)} rows={2} />

      <label>{t("form.relatedMaterial")}</label>
      <input value={relatedMaterialId} onChange={(e) => setRelatedMaterialId(e.target.value)} />

      <label>{t("form.usage")}</label>
      <div className="usage-options">
        {usageOptions.map((option) => (
          <label key={option} className="usage-option">
            <input
              type="radio"
              name="usage"
              checked={usage === option}
              onChange={() => setUsage(option)}
            />
            {t(`form.usageOptions.${option}`)}
          </label>
        ))}
      </div>

      <button type="submit" className="correction-submit" disabled={isSubmitting}>
        {isSubmitting ? t("form.submitting") : t("form.submit")}
      </button>
    </form>
  );
}
