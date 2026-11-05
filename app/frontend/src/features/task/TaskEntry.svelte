<script lang="ts">
  import { onMount } from "svelte";

  import AriadneForm from "../../components/AriadneForm.svelte";
  import TextField from "../../components/TextField.svelte";
  import TextareaField from "../../components/TextareaField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";
  import type { TaskValues } from "./TaskValues";

  import { maxLengthValidator } from "../../lib/validator";
  import type { Validity } from "../../lib/validity";
  import { PriorityDescriptionValidator } from "./PriorityDescriptionValidator.svelte";

  type Props = {
    submitLabel: string;
    initialValues: TaskValues;
    submitCondition?: (currentValues: TaskValues) => boolean;
    onConfirm?: () => Promise<void>;
  };

  let { submitLabel, initialValues, submitCondition, onConfirm }: Props = $props();

  // INFO: 相関Validotor
  const initialPriority = () => initialValues.priority;
  const initialDescription = () => initialValues.description;
  const priorityDescriptionValidator = new PriorityDescriptionValidator(initialPriority(), initialDescription());

  // INFO: 相関エラーをリアクティブに参照
  let correlationError = $derived.by(() => {
    const validation = priorityDescriptionValidator.getValidation();
    return validation.validity === "INVALID" ? (validation.errorMessage ?? null) : null;
  });

  // INFO: 相関Validityをリアクティブに参照
  let correlationValidity = $derived(priorityDescriptionValidator.getValidation().validity);

  // FUNCTION: onChangeイベントハンドラー
  function handlePriorityChange(value: number): void {
    priorityDescriptionValidator.setPriorityValue(value);
  }

  function handleDescriptionChange(value: string): void {
    priorityDescriptionValidator.setDescriptionValue(value);
  }

  // FUNCTION: onValidatedイベントハンドラー
  function handlePriorityValidated(validity: Validity): void {
    priorityDescriptionValidator.setPriorityValidity(validity);
  }

  function handleDescriptionValidated(validity: Validity): void {
    priorityDescriptionValidator.setDescriptionValidity(validity);
  }

  // FUNCTION: onFocusOutイベントハンドラー
  function handleCorrelationFocusOut(): void {
    priorityDescriptionValidator.validate();
  }

  // INFO: Field参照
  let titleField!: TextField;
  let descriptionField!: TextareaField;
  let priorityField!: TaskPriorityField;
  let statusField!: TaskStatusField;

  // INFO: リアクティブ参照の開始
  let mounted = $state(false);

  onMount(() => {
    mounted = true;
  });

  // INFO: Fieldの現在値をリアクティブに参照
  let currentValues = $derived(mounted ? getValues() : undefined);

  // INFO: FieldのValidityをリアクティブに参照
  let fieldValidities = $derived(
    mounted
      ? [
          titleField.getValidity(),
          descriptionField.getValidity(),
          priorityField.getValidity(),
          statusField.getValidity(),
        ]
      : [],
  );

  // INFO: Submit可否をリアクティブに判定
  let canSubmit = $derived(
    mounted &&
      fieldValidities.every((validity) => validity !== "INVALID") &&
      correlationValidity !== "INVALID" &&
      (currentValues !== undefined ? (submitCondition?.(currentValues) ?? true) : false),
  );

  // FUNCTION: 現在のTask入力値を集約
  export function getValues(): TaskValues {
    return {
      title: titleField.getValue(),
      description: descriptionField.getValue(),
      priority: priorityField.getValue(),
      status: statusField.getValue(),
    };
  }

  // FUNCTION: 全FieldのValidation実行
  function handleSubmit(): boolean {
    titleField.activateError();
    descriptionField.activateError();
    priorityField.activateError();
    statusField.activateError();

    titleField.validate();
    descriptionField.validate();
    priorityField.validate();
    statusField.validate();

    const allFieldsValid = [
      titleField.getValidity(),
      descriptionField.getValidity(),
      priorityField.getValidity(),
      statusField.getValidity(),
    ].every((validity) => validity === "VALID");

    if (!allFieldsValid) return false;

    priorityDescriptionValidator.activate();
    priorityDescriptionValidator.validate();

    return priorityDescriptionValidator.getValidation().validity === "VALID";
  }
</script>

<div class="entry-layout">
  <AriadneForm {submitLabel} {canSubmit} onSubmit={handleSubmit} {onConfirm}>
    <div class="entry-fields">
      <TextField
        bind:this={titleField}
        variant="edit"
        label="タイトル"
        required={true}
        value={initialValues.title}
        lines={2}
        maxLength={80}
        validators={[maxLengthValidator(50)]}
      />
      <TextareaField
        bind:this={descriptionField}
        variant="edit"
        label="説明"
        value={initialValues.description}
        lines={3}
        callbacks={{
          onChange: handleDescriptionChange,
          onValidated: handleDescriptionValidated,
          onFocusOut: handleCorrelationFocusOut,
        }}
        {correlationError}
      />
      <TaskPriorityField
        bind:this={priorityField}
        variant="edit"
        required={true}
        value={initialValues.priority}
        callbacks={{
          onChange: handlePriorityChange,
          onValidated: handlePriorityValidated,
          onFocusOut: handleCorrelationFocusOut,
        }}
        {correlationError}
      />
      <TaskStatusField bind:this={statusField} variant="edit" required={true} value={initialValues.status} />
    </div>
  </AriadneForm>
</div>

<style>
  .entry-layout {
    width: 100%;
    height: 100%;
    min-width: 0;
    min-height: 0;
  }

  .entry-fields {
    width: 500px;
    padding: 8px 6px;
  }
</style>
