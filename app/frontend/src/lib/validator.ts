import type { Validator } from "./validity";

// FUNCTION: 文字列必須チェック
export const requiredValidator: Validator<string> = (value) =>
  value.trim() === "" ? { validity: "INVALID", errorMessage: "必須項目です" } : { validity: "UNKNOWN" };
