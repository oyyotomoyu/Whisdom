import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { listModels } from "../../requests/system";
import type { ModelInfo } from "../../requests/system/types";
import { Spinner } from "../../components/Spinner";

export function SystemPanel() {
  const { t } = useTranslation("settings");
  const [models, setModels] = useState<ModelInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    listModels()
      .then(setModels)
      .finally(() => setIsLoading(false));
  }, []);

  if (isLoading) return <Spinner />;
  if (models.length === 0) return <p className="settings-empty">{t("system.empty")}</p>;

  return (
    <ul className="roles-list">
      {models.map((model) => (
        <li key={model.id} className="role-card">
          <span className="role-name">{model.name}</span>
          {model.active && <span className="active-badge">Active</span>}
        </li>
      ))}
    </ul>
  );
}
