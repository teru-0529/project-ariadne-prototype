<script lang="ts">
  import { onMount } from "svelte";

  import AriadneForm from "../../components/AriadneForm.svelte";
  import TextField from "../../components/TextField.svelte";
  import TextareaField from "../../components/TextareaField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";
  import type { TaskValues } from "./TaskValues";

  import { maxLengthValidator } from "../../lib/validator";

  type Props = {
    submitLabel: string;
    initialValues: TaskValues;
    submitCondition?: (currentValues: TaskValues) => boolean;
  };

  let { submitLabel, initialValues, submitCondition }: Props = $props();

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
</script>

<div class="entry-layout">
  <AriadneForm {submitLabel} {canSubmit}>
    <!-- TODO: ダミー -->
    <div class="callback-debug">
      <pre>fieldValidities: {fieldValidities.join(", ")}</pre>
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
      />
      <TaskPriorityField bind:this={priorityField} variant="edit" required={true} value={initialValues.priority} />
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
