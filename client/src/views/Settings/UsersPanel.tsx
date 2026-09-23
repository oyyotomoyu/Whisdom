import { useEffect, useState } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import { deleteUser, listUsers } from "../../requests/users";
import type { UserAccount } from "../../requests/users/types";
import { Spinner } from "../../components/Spinner";

export function UsersPanel() {
  const { t } = useTranslation("settings");
  const [users, setUsers] = useState<UserAccount[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    listUsers()
      .then(setUsers)
      .finally(() => setIsLoading(false));
  }, []);

  const handleDelete = async (id: string) => {
    setUsers((prev) => prev.filter((u) => u.id !== id));
    await deleteUser(id);
  };

  if (isLoading) return <Spinner />;
  if (users.length === 0) return <p className="settings-empty">{t("users.empty")}</p>;

  return (
    <table className="settings-table">
      <thead>
        <tr>
          <th>{t("users.table.name")}</th>
          <th>{t("users.table.email")}</th>
          <th>{t("users.table.role")}</th>
          <th>{t("users.table.status")}</th>
          <th>{t("users.table.actions")}</th>
        </tr>
      </thead>
      <tbody>
        {users.map((user) => (
          <tr key={user.id}>
            <td>{user.name}</td>
            <td>{user.email}</td>
            <td>{user.role}</td>
            <td>{user.active ? "Active" : "Disabled"}</td>
            <td>
              <button type="button" className="link-button" onClick={() => handleDelete(user.id)}>
                Delete
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
