import { useState, type FormEvent } from "react";
import { useTranslation } from "react-i18next";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import "./Login.css";

export function Login() {
  const { t } = useTranslation("login");
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const redirectTo = (location.state as { from?: Location })?.from?.pathname ?? "/chat";

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setIsSubmitting(true);
    setError(null);

    const result = await login(username, password);
    setIsSubmitting(false);

    if (result.ok) {
      navigate(redirectTo, { replace: true });
    } else {
      setError(t("invalidCredentials"));
    }
  };

  return (
    <form className="login-card" onSubmit={handleSubmit}>
      <h1 className="login-title">{t("title")}</h1>
      <p className="login-subtitle">{t("subtitle")}</p>

      <label className="login-label" htmlFor="username">
        {t("usernameLabel")}
      </label>
      <input
        id="username"
        name="username"
        type="text"
        autoComplete="username"
        className="login-input"
        placeholder={t("usernamePlaceholder")}
        value={username}
        onChange={(event) => setUsername(event.target.value)}
        required
      />

      <label className="login-label" htmlFor="password">
        {t("passwordLabel")}
      </label>
      <input
        id="password"
        name="password"
        type="password"
        autoComplete="current-password"
        className="login-input"
        placeholder={t("passwordPlaceholder")}
        value={password}
        onChange={(event) => setPassword(event.target.value)}
        required
      />

      {error && (
        <p className="login-error" role="alert">
          {error}
        </p>
      )}

      <button type="submit" className="login-submit" disabled={isSubmitting}>
        {isSubmitting ? t("submitting") : t("submit")}
      </button>
    </form>
  );
}
