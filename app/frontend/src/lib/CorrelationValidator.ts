import type { Validation } from "./validity";

// INFO: 相関Validotorの共通外部IF契約
export interface CorrelationValidator {
  activate(): void;
  validate(): void;
  getValidation(): Validation;
}

// INFO: 相関Validotorのフェーズ定義
export type CorrelationPhase = "INITIAL" | "ACTIVE";
