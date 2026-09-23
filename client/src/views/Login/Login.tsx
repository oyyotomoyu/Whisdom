import { useState, type FormEvent } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import { useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { getDevCredentials, isDevMode } from "../../requests/auth/devMode";
import "./Login.css";

export function Login() {
  const { t } = useTranslation("login");
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const devCredentials = getDevCredentials();
  const [username, setUsername] = useState(isDevMode() ? devCredentials.username : "");
  const [password, setPassword] = useState(isDevMode() ? devCredentials.password : "");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const redirectTo = (location.state as { from?: Location })?.from?.pathname ?? "/chat";
  const canUseDevLogin = isDevMode();

  const submitLogin = async (nextUsername: string, nextPassword: string) => {
    setIsSubmitting(true);
    setError(null);

    const result = await login(nextUsername, nextPassword);
    setIsSubmitting(false);

    if (result.ok) {
      navigate(redirectTo, { replace: true });
    } else {
      setError(t("invalidCredentials"));
    }
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    await submitLogin(username, password);
  };

  return (
    <section className="login-card" aria-labelledby="login-title">
      <div className="login-logo" aria-hidden="true">
        W
      </div>
      <div className="login-heading">
        <h1 className="login-title" id="login-title">
          {t("title")}
        </h1>
        <p className="login-subtitle">{t("subtitle")}</p>
      </div>

      <button
        type="button"
        className="login-sso"
        disabled={isSubmitting || !canUseDevLogin}
        onClick={() => submitLogin(devCredentials.username, devCredentials.password)}
      >
        {t("sso")}
      </button>

      <div className="login-divider">
        <span>{t("divider")}</span>
      </div>

      <form className="login-form" onSubmit={handleSubmit}>
        <label className="login-field" htmlFor="username">
          <span className="login-label">{t("usernameLabel")}</span>
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
        </label>

        <label className="login-field" htmlFor="password">
          <span className="login-label">{t("passwordLabel")}</span>
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
        </label>

        {error && (
          <p className="login-error" role="alert">
            {error}
          </p>
        )}

        <button type="submit" className="login-submit" disabled={isSubmitting}>
          {isSubmitting ? t("submitting") : t("submit")}
        </button>
      </form>
    </section>
  );
}
