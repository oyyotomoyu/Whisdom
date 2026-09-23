import { Outlet } from "react-router-dom";
import { useTranslation } from "../../../node_modules/react-i18next";
import { LanguageSelector } from "../../components/LanguageSelector";
import "./AuthLayout.css";

export function AuthLayout() {
  const { t } = useTranslation("global");

  return (
    <div className="auth-layout">
      <div className="auth-layout-top">
        <span className="app-name">{t("appName")}</span>
        <LanguageSelector className="app-language-selector" />
      </div>
      <div className="auth-layout-content">
        <Outlet />
      </div>
    </div>
  );
}
