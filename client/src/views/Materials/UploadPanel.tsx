import { useRef, useState, type DragEvent } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import { addMaterialSource, uploadMaterial, validateDestinationPath, validateSourceLocation } from "../../requests/materials";
import type { Material, MaterialSourceType } from "../../requests/materials/types";
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

const sourceTypes: MaterialSourceType[] = ["upload", "local_path", "wiki_url", "git_url", "url"];

export function UploadPanel({ onUploaded }: UploadPanelProps) {
  const { t } = useTranslation("materials");
  const inputRef = useRef<HTMLInputElement>(null);

  const [sourceType, setSourceType] = useState<MaterialSourceType>("upload");
  const [sourceLocation, setSourceLocation] = useState("");
  const [sourceError, setSourceError] = useState<string | null>(null);
  const [destinationPath, setDestinationPath] = useState("");
  const [destinationError, setDestinationError] = useState<string | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isSubmittingSource, setIsSubmittingSource] = useState(false);
  const [pending, setPending] = useState<PendingUpload[]>([]);

  const validateForm = () => {
    const pathError = validateDestinationPath(destinationPath);
    if (pathError) {
      setDestinationError(t(`upload.destinationError.${pathError}`));
      return false;
    }
    setDestinationError(null);

    const nextSourceError = validateSourceLocation(sourceType, sourceLocation);
    if (nextSourceError) {
      setSourceError(t(`upload.sourceError.${nextSourceError}`));
      return false;
    }
    setSourceError(null);
    return true;
  };

  const handleFiles = async (files: FileList | null) => {
    if (!files || files.length === 0 || !validateForm()) return;

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
          sourceType,
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

  const handleAddSource = async () => {
    if (sourceType === "upload" || !validateForm()) return;

    setIsSubmittingSource(true);
    try {
      const material = await addMaterialSource({
        sourceType,
        sourceLocation,
        destinationPath,
      });
      onUploaded(material);
      setSourceLocation("");
    } finally {
      setIsSubmittingSource(false);
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

      <label className="upload-field-label" htmlFor="sourceType">
        {t("upload.sourceTypeLabel")}
      </label>
      <select
        id="sourceType"
        className="upload-input"
        value={sourceType}
        onChange={(event) => {
          setSourceType(event.target.value as MaterialSourceType);
          setSourceError(null);
        }}
      >
        {sourceTypes.map((type) => (
          <option key={type} value={type}>
            {t(`upload.sourceTypes.${type}`)}
          </option>
        ))}
      </select>

      {sourceType !== "upload" && (
        <>
          <label className="upload-field-label" htmlFor="sourceLocation">
            {t("upload.sourceLocationLabel")}
          </label>
          <input
            id="sourceLocation"
            type="text"
            className="upload-input"
            placeholder={t(`upload.sourcePlaceholder.${sourceType}`)}
            value={sourceLocation}
            onChange={(event) => {
              setSourceLocation(event.target.value);
              setSourceError(null);
            }}
          />
          {sourceError && <p className="upload-error">{sourceError}</p>}
        </>
      )}

      <label className="upload-field-label" htmlFor="destinationPath">
        {t("upload.destinationLabel")}
      </label>
      <input
        id="destinationPath"
        type="text"
        className="upload-input"
        placeholder={t("upload.destinationPlaceholder")}
        value={destinationPath}
        onChange={(event) => {
          setDestinationPath(event.target.value);
          setDestinationError(null);
        }}
      />
      {destinationError && <p className="upload-error">{destinationError}</p>}

      {sourceType === "upload" ? (
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
      ) : (
        <button type="button" className="upload-submit" disabled={isSubmittingSource} onClick={handleAddSource}>
          {isSubmittingSource ? t("upload.addingSource") : t("upload.addSource")}
        </button>
      )}

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
