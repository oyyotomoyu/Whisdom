import { useEffect, useState } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import { listRoles } from "../../requests/roles";
import type { Role } from "../../requests/roles/types";
import { Spinner } from "../../components/Spinner";

export function RolesPanel() {
  const { t } = useTranslation("settings");
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    listRoles()
      .then(setRoles)
      .finally(() => setIsLoading(false));
  }, []);

  if (isLoading) return <Spinner />;
  if (roles.length === 0) return <p className="settings-empty">{t("roles.empty")}</p>;

  return (
    <ul className="roles-list">
      {roles.map((role) => (
        <li key={role.id} className="role-card">
          <span className="role-name">{role.name}</span>
          <span className="role-permissions">{role.permissions.join(", ")}</span>
        </li>
      ))}
    </ul>
  );
}
