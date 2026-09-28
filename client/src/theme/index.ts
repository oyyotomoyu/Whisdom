import theme from "@odm/theme.json";
import branding from "@odm/branding.json";
import "./tokens.css";

// ODM branding: odm/theme.json defines the active color palette so it can be
// swapped for a customer-specific build without touching component code or
// stylesheets. Applying it to the document root here, before anything
// renders, means every var(--color-*) reference in the app's CSS resolves to
// the ODM palette. See docs/client.md's "Theme And Branding (ODM)" section.
type ThemeColors = typeof theme.colors;

function toCssVariableName(token: string): string {
  return token.replace(/[A-Z]/g, (letter) => `-${letter.toLowerCase()}`);
}

for (const [token, value] of Object.entries(theme.colors)) {
  document.documentElement.style.setProperty(`--color-${toCssVariableName(token)}`, value);
}

document.documentElement.dataset.theme = theme.name;

// ODM branding also covers the product name: odm/branding.json, applied to
// the document title here so index.html never has to hardcode it. Views
// that need the name in visible text use t("appName", { productName }) (or
// an equivalent {{productName}} interpolation) rather than importing this
// directly, so the name stays translatable.
document.title = branding.productName;

export const productName = branding.productName;

export const breakpoints = {
  mobile: 768,
} as const;

/** Reads an ODM theme color for JS that can't use a CSS variable (canvas, charts, inline SVG fill). */
export function getThemeColor(token: keyof ThemeColors): string {
  return theme.colors[token];
}
