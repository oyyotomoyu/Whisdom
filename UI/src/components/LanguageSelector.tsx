import { useTranslation } from "react-i18next";
import { languages } from "../locales/languages";

export function LanguageSelector({ className }: { className?: string }) {
  const { i18n, t } = useTranslation("global");

  return (
    <select
      className={className}
      aria-label={t("language")}
      value={i18n.resolvedLanguage}
      onChange={(event) => i18n.changeLanguage(event.target.value)}
    >
      {languages.map((lang) => (
        <option key={lang.code} value={lang.code}>
          {lang.label}
        </option>
      ))}
    </select>
  );
}
