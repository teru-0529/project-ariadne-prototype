import type { Validation, Validity } from "../../lib/validity";
import type { CorrelationValidator, CorrelationPhase } from "../../lib/CorrelationValidator";

export class PriorityDescriptionValidator implements CorrelationValidator {
  // INFO: 相関Validationの活性状態
  // private phase: CorrelationPhase = "INITIAL";
  // TODO:一時的にリアクティブ
  private phase: CorrelationPhase = $state("INITIAL");

  // INFO: 相関Validationの判定結果
  private validation = $state<Validation>({
    validity: "UNKNOWN",
  });

  // INFO: 相関Validationの対象Field
  private priority: number;
  private priorityValidity: Validity = "UNKNOWN";

  private description: string;
  private descriptionValidity: Validity = "UNKNOWN";

  // FUNCTION: コンストラクタ（対象Fieldの初期値を設定する）
  constructor(priority: number, description: string) {
    this.priority = priority;
    this.description = description;
  }

  // FUNCTION: Priorityの値・Validity更新
  setPriorityValue(value: number): void {
    this.priority = value;
  }

  setPriorityValidity(validity: Validity): void {
    this.priorityValidity = validity;
  }

  // FUNCTION: Descriptionの値・Validity更新
  setDescriptionValue(value: string): void {
    this.description = value;
  }

  setDescriptionValidity(validity: Validity): void {
    this.descriptionValidity = validity;
  }

  // FUNCTION: 相関Validationの活性化
  activate(): void {
    this.phase = "ACTIVE";
  }

  // FUNCTION: 相関Validationの評価
  validate(): void {
    // TODO: A07-11
  }

  // FUNCTION: 相関Validationの判定結果取得
  getValidation(): Validation {
    return this.validation;
  }

  // TODO: A07-18 初期状態確認用
  getPhase(): CorrelationPhase {
    return this.phase;
  }
}
