import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { Spinner } from "../components/Spinner";

export function AdminRoute() {
  const { isAdmin, isAuthLoading } = useAuth();

  if (isAuthLoading) {
    return <Spinner fullscreen />;
  }

  if (!isAdmin) {
    return <Navigate to="/unauthorized" replace />;
  }

  return <Outlet />;
}
