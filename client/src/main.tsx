import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import faviconUrl from "@odm/img/favicon.svg?url";
import "./index.css";
import "./theme/index.ts";
import "./locales/i18n.ts";
import App from "./App.tsx";

// The favicon is an ODM asset (see docs/client.md), so it's injected here
// rather than declared as a static <link> in index.html: that way an ODM
// build only has to replace odm/img/favicon.svg, never index.html.
const favicon = document.createElement("link");
favicon.rel = "icon";
favicon.type = "image/svg+xml";
favicon.href = faviconUrl;
document.head.appendChild(favicon);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
