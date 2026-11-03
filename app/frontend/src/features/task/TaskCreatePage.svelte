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

  // FUNCTION: Page離脱
  export function leave(callback: () => void): void {
    // TODO: 検証用
    window.alert("TaskCreatePage.leave()");
    callback();
  }
</script>

<AriadnePage title="Task 作成">
  <TaskEntry submitLabel="Create Task" {initialValues} />
</AriadnePage>
