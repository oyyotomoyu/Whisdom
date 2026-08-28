import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useSearchParams } from "react-router-dom";
import { createCorrection, deleteCorrection, listCorrections } from "../../requests/corrections";
import type { Correction } from "../../requests/corrections/types";
import { CorrectionForm } from "./CorrectionForm";
import { CorrectionsList } from "./CorrectionsList";
import { Spinner } from "../../components/Spinner";
import "./Corrections.css";

export function Corrections() {
  const { t } = useTranslation("corrections");
  const [searchParams] = useSearchParams();
  const [corrections, setCorrections] = useState<Correction[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const prefillAnswer = searchParams.get("answer") ?? undefined;
  const prefillMessageId = searchParams.get("messageId") ?? undefined;

  useEffect(() => {
    listCorrections()
      .then(setCorrections)
      .finally(() => setIsLoading(false));
  }, []);

  const handleCreate = async (values: Parameters<typeof createCorrection>[0]) => {
    const created = await createCorrection(values);
    setCorrections((prev) => [created, ...prev]);
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm(t("deleteConfirm"))) return;
    setCorrections((prev) => prev.filter((c) => c.id !== id));
    await deleteCorrection(id);
  };

  return (
    <div className="corrections-view">
      <header className="corrections-header">
        <h1>{t("title")}</h1>
        <p>{t("subtitle")}</p>
      </header>

      <CorrectionForm initialAnswer={prefillAnswer} initialMessageId={prefillMessageId} onSubmit={handleCreate} />

      {isLoading ? <Spinner /> : <CorrectionsList corrections={corrections} onDelete={handleDelete} />}
    </div>
  );
}
