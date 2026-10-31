// INFO: 検証状態
export type Validity = "UNKNOWN" | "VALID" | "INVALID";

// INFO: 検証結果（INVALIDの場合のみエラーメッセージを保持）
export type Validation =
  { validity: "UNKNOWN" } | { validity: "VALID" } | { validity: "INVALID"; errorMessage: string };

// INFO: Validation関数の型
export type Validator<T> = (value: T) => Validation;
