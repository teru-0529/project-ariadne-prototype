<script lang="ts">
  import AriadneForm from "../../components/AriadneForm.svelte";
  import TextField from "../../components/TextField.svelte";
  import TextareaField from "../../components/TextareaField.svelte";
  import TaskPriorityField from "./TaskPriorityField.svelte";
  import TaskStatusField from "./TaskStatusField.svelte";
</script>

<AriadneForm submitLabel="Create Task">
  <TextField variant="edit" label="タイトル" value="" />
  <TextareaField variant="edit" label="説明" value="" />
  <TaskPriorityField variant="edit" label="優先度" value={3} />
  <!--TODO: Field群は次の増分で追加 -->
</AriadneForm>
