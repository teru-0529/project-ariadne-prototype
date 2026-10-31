<script lang="ts">
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
  };

  let { submitLabel, initialValues }: Props = $props();

  // INFO: Field参照
  let titleField: TextField;
  let descriptionField: TextareaField;
  let priorityField: TaskPriorityField;
  let statusField: TaskStatusField;

  // FUNCTION: 現在のTask入力値を集約
  function getValues(): TaskValues {
    return {
      title: titleField.getValue(),
      description: descriptionField.getValue(),
      priority: priorityField.getValue(),
      status: statusField.getValue(),
    };
  }

  // TODO: A07-12 動作確認用
  let debugValues = $state<TaskValues | null>(null);

  // TODO: A07-12 動作確認用
  function handleDebugGetValues() {
    debugValues = getValues();
  }
</script>

<div class="entry-layout">
  <AriadneForm {submitLabel}>
    <!-- TODO: A07-12 動作確認用 -->
    <div class="callback-debug">
      <button type="button" onclick={handleDebugGetValues}> getValues確認 </button>

      {#if debugValues}
        <pre>{JSON.stringify(debugValues, null, 2)}</pre>
      {/if}
    </div>
    <!-- TODO: A07-12 動作確認用 -->

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
