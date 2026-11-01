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
  };

  let { submitLabel, initialValues, submitCondition }: Props = $props();

  // INFO: 相関Validotor
  const initialPriority = () => initialValues.priority;
  const initialDescription = () => initialValues.description;
  const priorityDescriptionValidator = new PriorityDescriptionValidator(initialPriority(), initialDescription());

  // FUNCTION: onChangeイベントハンドラー
  function handlePriorityChange(value: number): void {
    correlationDebugVersion += 1;
    priorityDescriptionValidator.setPriorityValue(value);
  }

  function handleDescriptionChange(value: string): void {
    correlationDebugVersion += 1;
    priorityDescriptionValidator.setDescriptionValue(value);
  }

  // FUNCTION: onValidatedイベントハンドラー
  function handlePriorityValidated(validity: Validity): void {
    correlationDebugVersion += 1;
    priorityDescriptionValidator.setPriorityValidity(validity);
  }

  function handleDescriptionValidated(validity: Validity): void {
    correlationDebugVersion += 1;
    priorityDescriptionValidator.setDescriptionValidity(validity);
  }

  // TODO: A07-09 動作確認用
  let correlationDebugVersion = $state(0);

  // TODO: A07-09 動作確認用
  let correlationDebugState = $derived.by(() => {
    const version = correlationDebugVersion;
    return {
      version,
      ...priorityDescriptionValidator.getDebugState(),
    };
  });

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

  // INFO: Submit可否をリアクティブに判定 TODO: 相関ValidatorのValidityを追加
  let canSubmit = $derived(
    mounted &&
      fieldValidities.every((validity) => validity !== "INVALID") &&
      (currentValues !== undefined ? (submitCondition?.(currentValues) ?? true) : false),
  );

  // FUNCTION: 現在のTask入力値を集約
  function getValues(): TaskValues {
    return {
      title: titleField.getValue(),
      description: descriptionField.getValue(),
      priority: priorityField.getValue(),
      status: statusField.getValue(),
    };
  }

  // FUNCTION: 未保存変更の有無
  export function hasUnsavedChanges(): boolean {
    const currentValues = getValues();

    return (
      currentValues.title !== initialValues.title ||
      currentValues.description !== initialValues.description ||
      currentValues.priority !== initialValues.priority ||
      currentValues.status !== initialValues.status
    );
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

    return true;
  }
</script>

<div class="entry-layout">
  <AriadneForm {submitLabel} {canSubmit} onSubmit={handleSubmit}>
    <!-- TODO: ダミー -->
    <div class="callback-debug">
      <pre>fieldValidities: {fieldValidities.join(", ")}</pre>
      <pre>correlationValidity: {priorityDescriptionValidator.getValidation().validity}</pre>
      <pre>correlationDebug: {JSON.stringify(correlationDebugState)}</pre>
    </div>
    <!-- TODO: ダミー -->

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
        callbacks={{ onChange: handleDescriptionChange, onValidated: handleDescriptionValidated }}
      />
      <TaskPriorityField
        bind:this={priorityField}
        variant="edit"
        required={true}
        value={initialValues.priority}
        callbacks={{ onChange: handlePriorityChange, onValidated: handlePriorityValidated }}
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

  /* TODO: debug */
  .callback-debug {
    padding: 8px 6px;
    font-size: 10px;
    line-height: 1.8;
  }
</style>
