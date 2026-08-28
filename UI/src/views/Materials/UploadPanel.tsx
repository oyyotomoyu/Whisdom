import { useRef, useState, type DragEvent } from "react";
import { useTranslation } from "react-i18next";
import { uploadMaterial, validateDestinationPath } from "../../requests/materials";
import type { Material } from "../../requests/materials/types";
import "./UploadPanel.css";

interface UploadPanelProps {
  onUploaded: (material: Material) => void;
}

interface PendingUpload {
  id: string;
  file: File;
  progress: number;
  error?: string;
}

export function UploadPanel({ onUploaded }: UploadPanelProps) {
  const { t } = useTranslation("materials");
  const inputRef = useRef<HTMLInputElement>(null);

  const [destinationPath, setDestinationPath] = useState("");
  const [destinationError, setDestinationError] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [pending, setPending] = useState<PendingUpload[]>([]);

  const handleFiles = async (files: FileList | null) => {
    if (!files || files.length === 0) return;

    const pathError = validateDestinationPath(destinationPath);
    if (pathError) {
      setDestinationError(t(`upload.destinationError.${pathError}`));
      return;
    }
    setDestinationError(null);

    const uploads: PendingUpload[] = Array.from(files).map((file) => ({
      id: `${file.name}_${file.size}_${Date.now()}`,
      file,
      progress: 0,
    }));
    setPending((prev) => [...prev, ...uploads]);

    for (const upload of uploads) {
      try {
        const material = await uploadMaterial({
          file: upload.file,
          destinationPath,
          onProgress: (percent) =>
            setPending((prev) => prev.map((p) => (p.id === upload.id ? { ...p, progress: percent } : p))),
        });
        onUploaded(material);
        setPending((prev) => prev.filter((p) => p.id !== upload.id));
      } catch {
        setPending((prev) => prev.map((p) => (p.id === upload.id ? { ...p, error: "global:errors.generic" } : p)));
      }
    }
  };

  const handleDrop = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setIsDragging(false);
    handleFiles(event.dataTransfer.files);
  };

  return (
    <div className="upload-panel">
      <h2 className="upload-panel-title">{t("upload.title")}</h2>

      <label className="upload-field-label" htmlFor="destinationPath">
        {t("upload.destinationLabel")}
      </label>
      <input
        id="destinationPath"
        type="text"
        className="upload-destination-input"
        placeholder={t("upload.destinationPlaceholder")}
        value={destinationPath}
        onChange={(event) => {
          setDestinationPath(event.target.value);
          setDestinationError(null);
        }}
      />
      {destinationError && <p className="upload-error">{destinationError}</p>}

      <div
        className={`upload-dropzone${isDragging ? " dragging" : ""}`}
        onDragOver={(event) => {
          event.preventDefault();
          setIsDragging(true);
        }}
        onDragLeave={() => setIsDragging(false)}
        onDrop={handleDrop}
        onClick={() => inputRef.current?.click()}
        role="button"
        tabIndex={0}
      >
        <p>{t("upload.dragHint")}</p>
        <span className="upload-browse-link">{t("upload.browse")}</span>
        <input
          ref={inputRef}
          type="file"
          multiple
          hidden
          onChange={(event) => {
            handleFiles(event.target.files);
            event.target.value = "";
          }}
        />
      </div>

      {pending.length > 0 && (
        <ul className="upload-pending-list">
          {pending.map((upload) => (
            <li key={upload.id} className={upload.error ? "error" : ""}>
              <span className="upload-pending-name">{upload.file.name}</span>
              {upload.error ? (
                <span className="upload-pending-error">{t(upload.error)}</span>
              ) : (
                <div className="upload-progress-track">
                  <div className="upload-progress-fill" style={{ width: `${upload.progress}%` }} />
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
