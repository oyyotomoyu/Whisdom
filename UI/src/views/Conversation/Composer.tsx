import { useState, type KeyboardEvent } from "react";
import { useTranslation } from "react-i18next";
import "./Composer.css";

interface ComposerProps {
  onSend: (text: string) => void;
  isSending: boolean;
}

export function Composer({ onSend, isSending }: ComposerProps) {
  const { t } = useTranslation("conversation");
  const [value, setValue] = useState("");

  const handleSend = () => {
    if (!value.trim() || isSending) return;
    onSend(value);
    setValue("");
  };

  const handleKeyDown = (event: KeyboardEvent<HTMLTextAreaElement>) => {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="composer">
      <div className="composer-inner">
        <textarea
          className="composer-input"
          placeholder={t("composer.placeholder")}
          value={value}
          onChange={(event) => setValue(event.target.value)}
          onKeyDown={handleKeyDown}
          disabled={isSending}
          rows={1}
        />
        <button
          type="button"
          className="composer-send"
          onClick={handleSend}
          disabled={isSending || !value.trim()}
          aria-label={t("composer.send")}
        >
          {t("composer.send")}
        </button>
      </div>
      <p className="composer-hint">{t("composer.hint")}</p>
    </div>
  );
}
