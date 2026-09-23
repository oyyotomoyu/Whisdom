import { createBrowserRouter, Navigate, RouterProvider } from "react-router-dom";
import { AuthLayout } from "../layouts/AuthLayout";
import { AppLayout } from "../layouts/AppLayout";
import { ProtectedRoute } from "../routes/ProtectedRoute";
import { AdminRoute } from "../routes/AdminRoute";
import { Login } from "./Login";
import { Conversation } from "./Conversation";
import { Materials } from "./Materials";
import { Corrections } from "./Corrections";
import { Settings } from "./Settings";
import { Unauthorized } from "./Unauthorized";
import { NotFound } from "./NotFound";

const router = createBrowserRouter([
  {
    element: <AuthLayout />,
    children: [{ path: "/login", element: <Login /> }],
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <AppLayout />,
        children: [
          { path: "/", element: <Navigate to="/chat" replace /> },
          { path: "/chat", element: <Conversation /> },
          { path: "/chat/:conversationId", element: <Conversation /> },
          { path: "/unauthorized", element: <Unauthorized /> },
          {
            element: <AdminRoute />,
            children: [
              { path: "/admin/materials", element: <Materials /> },
              { path: "/admin/corrections", element: <Corrections /> },
              { path: "/settings", element: <Settings /> },
            ],
          },
        ],
      },
    ],
  },
  { path: "*", element: <NotFound /> },
]);

export function AppRoutes() {
  return <RouterProvider router={router} />;
}
