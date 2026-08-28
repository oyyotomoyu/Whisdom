import { useState } from "react";
import { useTranslation } from "react-i18next";
import type { Material } from "../../requests/materials/types";
import { validateDestinationPath } from "../../requests/materials";
import "./MaterialsTable.css";

interface MaterialsTableProps {
  materials: Material[];
  onDelete: (id: string) => void;
  onReprocess: (id: string) => void;
  onUpdateDestination: (id: string, destinationPath: string) => Promise<void>;
}

export function MaterialsTable({ materials, onDelete, onReprocess, onUpdateDestination }: MaterialsTableProps) {
  const { t } = useTranslation("materials");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [draftPath, setDraftPath] = useState("");
  const [editError, setEditError] = useState<string | null>(null);

  const startEdit = (material: Material) => {
    setEditingId(material.id);
    setDraftPath(material.destinationPath);
    setEditError(null);
  };

  const saveEdit = async (id: string) => {
    const pathError = validateDestinationPath(draftPath);
    if (pathError) {
      setEditError(t(`upload.destinationError.${pathError}`));
      return;
    }
    await onUpdateDestination(id, draftPath);
    setEditingId(null);
  };

  if (materials.length === 0) {
    return <p className="materials-table-empty">{t("table.empty")}</p>;
  }

  return (
    <div className="materials-table-wrap">
      <table className="materials-table">
        <thead>
          <tr>
            <th>{t("table.name")}</th>
            <th>{t("table.type")}</th>
            <th>{t("table.uploadedAt")}</th>
            <th>{t("table.uploadedBy")}</th>
            <th>{t("table.status")}</th>
            <th>{t("table.destination")}</th>
            <th>{t("table.rag")}</th>
            <th>{t("table.training")}</th>
            <th>{t("table.actions")}</th>
          </tr>
        </thead>
        <tbody>
          {materials.map((material) => (
            <tr key={material.id}>
              <td>{material.filename}</td>
              <td>{material.type}</td>
              <td>{new Date(material.uploadedAt).toLocaleString()}</td>
              <td>{material.uploadedBy}</td>
              <td>
                <span className={`status-pill status-${material.status}`}>{t(`status.${material.status}`)}</span>
              </td>
              <td>
                {editingId === material.id ? (
                  <div className="destination-editor">
                    <input
                      value={draftPath}
                      onChange={(event) => setDraftPath(event.target.value)}
                      autoFocus
                    />
                    {editError && <p className="upload-error">{editError}</p>}
                    <div className="destination-editor-actions">
                      <button type="button" onClick={() => saveEdit(material.id)}>
                        Save
                      </button>
                      <button type="button" onClick={() => setEditingId(null)}>
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  <button type="button" className="link-button" onClick={() => startEdit(material)}>
                    {material.destinationPath}
                  </button>
                )}
              </td>
              <td>{material.ragAvailable ? "✓" : "—"}</td>
              <td>{material.trainingAvailable ? "✓" : "—"}</td>
              <td>
                <div className="materials-row-actions">
                  <button type="button" onClick={() => onReprocess(material.id)}>
                    {t("actions.reprocess")}
                  </button>
                  <button
                    type="button"
                    className="danger"
                    onClick={() => {
                      if (window.confirm(t("deleteConfirm"))) onDelete(material.id);
                    }}
                  >
                    {t("actions.delete")}
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
