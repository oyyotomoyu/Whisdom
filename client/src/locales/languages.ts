export interface LanguageOption {
  code: string;
  label: string;
}

export const languages: LanguageOption[] = [
  { code: "en", label: "English" },
  { code: "zh-TW", label: "繁體中文" },
  { code: "ja", label: "日本語" },
];

export const defaultLanguage = "en";
