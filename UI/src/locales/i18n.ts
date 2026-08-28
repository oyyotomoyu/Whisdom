import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import globalEn from "./key/Global/en";
import globalZhTw from "./key/Global/zh-TW";
import globalJa from "./key/Global/ja";

import loginEn from "./key/Login/en";
import loginZhTw from "./key/Login/zh-TW";
import loginJa from "./key/Login/ja";

import conversationEn from "./key/Conversation/en";
import conversationZhTw from "./key/Conversation/zh-TW";
import conversationJa from "./key/Conversation/ja";

import adminEn from "./key/Admin/en";
import adminZhTw from "./key/Admin/zh-TW";
import adminJa from "./key/Admin/ja";

import materialsEn from "./key/Materials/en";
import materialsZhTw from "./key/Materials/zh-TW";
import materialsJa from "./key/Materials/ja";

import correctionsEn from "./key/Corrections/en";
import correctionsZhTw from "./key/Corrections/zh-TW";
import correctionsJa from "./key/Corrections/ja";

import settingsEn from "./key/Settings/en";
import settingsZhTw from "./key/Settings/zh-TW";
import settingsJa from "./key/Settings/ja";

import { defaultLanguage } from "./languages";

export const defaultNS = "global";

export const resources = {
  en: {
    global: globalEn,
    login: loginEn,
    conversation: conversationEn,
    admin: adminEn,
    materials: materialsEn,
    corrections: correctionsEn,
    settings: settingsEn,
  },
  "zh-TW": {
    global: globalZhTw,
    login: loginZhTw,
    conversation: conversationZhTw,
    admin: adminZhTw,
    materials: materialsZhTw,
    corrections: correctionsZhTw,
    settings: settingsZhTw,
  },
  ja: {
    global: globalJa,
    login: loginJa,
    conversation: conversationJa,
    admin: adminJa,
    materials: materialsJa,
    corrections: correctionsJa,
    settings: settingsJa,
  },
} as const;

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: defaultLanguage,
    defaultNS,
    ns: ["global", "login", "conversation", "admin", "materials", "corrections", "settings"],
    interpolation: {
      escapeValue: false,
    },
    detection: {
      order: ["localStorage", "navigator"],
      caches: ["localStorage"],
      lookupLocalStorage: "whisdom_language",
    },
  });

export default i18n;
