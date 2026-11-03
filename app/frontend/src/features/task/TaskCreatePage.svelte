<script lang="ts">
  import AriadnePage from "../../components/AriadnePage.svelte";
  import TaskEntry from "./TaskEntry.svelte";
  import type { TaskValues } from "./TaskValues";

  // INFO: 登録時初期値
  const initialValues: TaskValues = {
    title: "",
    description: "",
    priority: 3,
    status: "TODO",
  };

  // INFO: taskEntryへの参照
  let taskEntry = $state<TaskEntry>();

  // FUNCTION: Page離脱
  export function leave(callback: () => void): void {
    // INFO: 未保存変更の有無を確認
    const changed = hasUnsavedChanges();

    // TODO: 検証用（A12-04でModalに置き換える）
    window.alert(`TaskCreatePage.leave(): changed=${changed}`);
    callback();
  }

  // FUNCTION: 未保存変更の有無
  function hasUnsavedChanges(): boolean {
    const currentValues = taskEntry!.getValues();

    return (
      currentValues.title !== initialValues.title ||
      currentValues.description !== initialValues.description ||
      currentValues.priority !== initialValues.priority ||
      currentValues.status !== initialValues.status
    );
  }
</script>

<AriadnePage title="Task 作成">
  <TaskEntry bind:this={taskEntry} submitLabel="Create Task" {initialValues} />
</AriadnePage>
