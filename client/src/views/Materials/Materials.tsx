import { useEffect, useState } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import {
  deleteMaterial,
  listMaterials,
  reprocessMaterial,
  updateMaterialDestination,
} from "../../requests/materials";
import type { Material } from "../../requests/materials/types";
import { UploadPanel } from "./UploadPanel";
import { MaterialsTable } from "./MaterialsTable";
import { Spinner } from "../../components/Spinner";
import "./Materials.css";

export function Materials() {
  const { t } = useTranslation("materials");
  const [materials, setMaterials] = useState<Material[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    listMaterials()
      .then(setMaterials)
      .finally(() => setIsLoading(false));
  }, []);

  const handleUploaded = (material: Material) => {
    setMaterials((prev) => [material, ...prev]);
  };

  const handleDelete = async (id: string) => {
    setMaterials((prev) => prev.filter((m) => m.id !== id));
    await deleteMaterial(id);
  };

  const handleReprocess = async (id: string) => {
    setMaterials((prev) => prev.map((m) => (m.id === id ? { ...m, status: "processing" } : m)));
    await reprocessMaterial(id);
  };

  const handleUpdateDestination = async (id: string, destinationPath: string) => {
    const updated = await updateMaterialDestination(id, destinationPath);
    setMaterials((prev) => prev.map((m) => (m.id === id ? updated : m)));
  };

  return (
    <div className="materials-view">
      <header className="materials-header">
        <h1>{t("title")}</h1>
        <p>{t("subtitle")}</p>
      </header>

      <UploadPanel onUploaded={handleUploaded} />

      {isLoading ? (
        <Spinner />
      ) : (
        <MaterialsTable
          materials={materials}
          onDelete={handleDelete}
          onReprocess={handleReprocess}
          onUpdateDestination={handleUpdateDestination}
        />
      )}
    </div>
  );
}
