import { NavLink, Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useAuth } from "../../hooks/useAuth";
import { useIsMobile } from "../../hooks/useMediaQuery";
import { useUiStore } from "../../store/uiStore";
import { ConversationSidebar } from "../../components/ConversationSidebar";
import { LanguageSelector } from "../../components/LanguageSelector";
import "./AppLayout.css";

export function AppLayout() {
  const { t } = useTranslation("global");
  const { user, isAdmin, logout } = useAuth();
  const isMobile = useIsMobile();
  const isSidebarOpen = useUiStore((s) => s.isSidebarOpen);
  const closeSidebar = useUiStore((s) => s.closeSidebar);
  const toggleSidebar = useUiStore((s) => s.toggleSidebar);

  const sidebarContent = (
    <>
      <div className="app-sidebar-header">
        <span className="app-name">{t("appName")}</span>
      </div>
      <div className="app-sidebar-conversations">
        <ConversationSidebar onNavigate={isMobile ? closeSidebar : undefined} />
      </div>
      {isAdmin && (
        <div className="app-sidebar-admin">
          <div className="sidebar-section-label">{t("nav.admin")}</div>
          <NavLink to="/admin/materials" className="app-nav-link" onClick={isMobile ? closeSidebar : undefined}>
            {t("nav.materials")}
          </NavLink>
          <NavLink to="/admin/corrections" className="app-nav-link" onClick={isMobile ? closeSidebar : undefined}>
            {t("nav.corrections")}
          </NavLink>
          <NavLink to="/settings" className="app-nav-link" onClick={isMobile ? closeSidebar : undefined}>
            {t("nav.settings")}
          </NavLink>
        </div>
      )}
      <div className="app-sidebar-footer">
        <LanguageSelector className="app-language-selector" />
        <div className="app-user">
          <span className="app-user-name">{user?.name}</span>
          <button type="button" className="app-logout" onClick={logout}>
            {t("nav.logout")}
          </button>
        </div>
      </div>
    </>
  );

  return (
    <div className="app-layout">
      {isMobile && (
        <header className="app-topbar">
          <button type="button" className="app-topbar-menu" aria-label="Open menu" onClick={toggleSidebar}>
            ☰
          </button>
          <span className="app-name">{t("appName")}</span>
        </header>
      )}

      {isMobile ? (
        <>
          <div className={`app-drawer${isSidebarOpen ? " open" : ""}`}>{sidebarContent}</div>
          {isSidebarOpen && <div className="app-drawer-backdrop" onClick={closeSidebar} />}
        </>
      ) : (
        <aside className="app-sidebar">{sidebarContent}</aside>
      )}

      <main className="app-main">
        <Outlet />
      </main>
    </div>
  );
}
