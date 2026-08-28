import { useState } from "react";
import { useTranslation } from "react-i18next";
import { UsersPanel } from "./UsersPanel";
import { RolesPanel } from "./RolesPanel";
import { SystemPanel } from "./SystemPanel";
import "./Settings.css";

type Tab = "users" | "roles" | "system";

export function Settings() {
  const { t } = useTranslation("settings");
  const [tab, setTab] = useState<Tab>("users");

  return (
    <div className="settings-view">
      <header className="settings-header">
        <h1>{t("title")}</h1>
      </header>

      <div className="settings-tabs">
        {(["users", "roles", "system"] as Tab[]).map((key) => (
          <button
            key={key}
            type="button"
            className={`settings-tab${tab === key ? " active" : ""}`}
            onClick={() => setTab(key)}
          >
            {t(`tabs.${key}`)}
          </button>
        ))}
      </div>

      <div className="settings-panel">
        {tab === "users" && <UsersPanel />}
        {tab === "roles" && <RolesPanel />}
        {tab === "system" && <SystemPanel />}
      </div>
    </div>
  );
}
