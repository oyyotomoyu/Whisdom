/// <reference types="vite/client" />

declare const process: {
  env: {
    NODE_ENV?: string;
  };
};

// Vite resolves a "?url" suffixed import to the asset's built URL at
// runtime, but ships no ambient type for it (see vite/client.d.ts, which
// only covers bare extensions). Used for ODM branding assets imported as
// e.g. `@odm/img/favicon.svg?url` (see client/src/main.tsx).
declare module "*?url" {
  const src: string;
  export default src;
}
