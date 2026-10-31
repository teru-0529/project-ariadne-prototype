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
</script>

<div class="entry-layout">
  <AriadneForm {submitLabel}>
    <div class="entry-fields">
      <TextField
        variant="edit"
        label="タイトル"
        required={true}
        value={initialValues.title}
        lines={2}
        maxLength={80}
        validators={[maxLengthValidator(50)]}
      />
      <TextareaField variant="edit" label="説明" value={initialValues.description} lines={3} />
      <TaskPriorityField variant="edit" required={true} value={initialValues.priority} />
      <TaskStatusField variant="edit" required={true} value={initialValues.status} />
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
