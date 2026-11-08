import type { Validity } from "./validity";

// INFO: Fieldから外部への通知契約
export type FieldCallbacks<T> = {
  onChange?: (value: T) => void;
  onValidated?: (validity: Validity) => void;
  onFocusOut?: () => void;
};
