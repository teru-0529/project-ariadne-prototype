import type { Validation, Validator } from "./validity";

// FUNCTION: 文字列必須チェック
export const requiredValidator: Validator<string> = (value) =>
  value.trim() === "" ? { validity: "INVALID", errorMessage: "必須項目です" } : { validity: "UNKNOWN" };

// FUNCTION: 最大文字数チェック
export const maxLengthValidator =
  (maxLength: number): Validator<string> =>
  (value) =>
    value.length > maxLength
      ? { validity: "INVALID", errorMessage: `${maxLength}文字以内で入力してください` }
      : { validity: "UNKNOWN" };

// FUNCTION: Validatorの順次実行
export function runValidators<T>(value: T, validators: Validator<T>[]): Validation {
  for (const validator of validators) {
    const result = validator(value);

    if (result.validity === "INVALID") {
      return result;
    }
  }

  return { validity: "VALID" };
}
