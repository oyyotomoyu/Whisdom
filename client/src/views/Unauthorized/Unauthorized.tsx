import { Link } from "react-router-dom";
import { useTranslation } from "../../../node_modules/react-i18next";
import { EmptyState } from "../../components/EmptyState";

export function Unauthorized() {
  const { t } = useTranslation("admin");

  return (
    <div style={{ height: "100%" }}>
      <EmptyState
        title={t("unauthorized.title")}
        subtitle={t("unauthorized.subtitle")}
        action={<Link to="/chat">{t("unauthorized.backToChat")}</Link>}
      />
    </div>
  );
}
