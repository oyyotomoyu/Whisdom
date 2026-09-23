import { useEffect } from "react";
import { AppRoutes } from "./views";
import { useAuthStore } from "./store/authStore";
import { fetchCurrentUser } from "./requests/auth";
import { clearAccessToken, setAccessToken } from "./requests/core";

function App() {
  const accessToken = useAuthStore((state) => state.accessToken);
  const setSession = useAuthStore((state) => state.setSession);
  const clearSession = useAuthStore((state) => state.clearSession);
  const setAuthLoading = useAuthStore((state) => state.setAuthLoading);

  useEffect(() => {
    if (!accessToken) {
      clearAccessToken();
      setAuthLoading(false);
      return;
    }

    setAccessToken(accessToken);
    fetchCurrentUser()
      .then((user) => setSession(user, accessToken))
      .catch(() => {
        clearAccessToken();
        clearSession();
      });
    // Runs once on mount to validate any persisted session.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return <AppRoutes />;
}

export default App;
