<script lang="ts">
  import SideNavigation from "../../components/SideNavigation.svelte";
  import TaskListPage from "./TaskListPage.svelte";

  import { TaskScreen, TaskContext } from "./TaskContext";
  import { toErrorMessage } from "../../lib/errors/ProcessedError";

  export let onError: (message: string) => void;

  // INFO: Side Navigation に渡すリスト
  const navigationItems: { value: TaskScreen; label: string }[] = [
    { value: TaskScreen.List, label: "List" },
    { value: TaskScreen.Create, label: "Create" },
  ];

  // INFO: 画面遷移中
  let isNavigating = false;

  // INFO: 選択中の画面
  let currentTaskOperation: TaskScreen | null = null;

  // INFO: タスクコンテキスト
  let context = new TaskContext();

  // FUNCTION: 画面遷移
  async function navigate(value: string): Promise<void> {
    // INFO: 二重操作防止
    if (isNavigating) {
      return;
    }

    isNavigating = true;

    try {
      // INFO: 選択された画面の検証
      if (!navigationItems.some((item) => item.value === value)) {
        console.error("TaskManagement: Invalid screen", value);
        throw new Error("予期しないエラーが発生しました。");
      }

      const screen = value as TaskScreen;

      await context.prepare(screen);
      context.commit(screen);
      currentTaskOperation = screen;
    } catch (error) {
      onError(toErrorMessage(error));
    } finally {
      isNavigating = false;
    }
  }
</script>

<div class="feature-area with-side-navigation">
  <SideNavigation label="Task" items={navigationItems} current={currentTaskOperation} onSelect={navigate} />

  <main class="main-area">
    {#if currentTaskOperation === TaskScreen.List}
      <TaskListPage {context} />
    {:else if currentTaskOperation === TaskScreen.Create}
      <!-- TODO: Createのダミー実装  -->
      <p>Create Task.</p>
    {:else}
      <p>Select an operation.</p>
    {/if}
  </main>
</div>

<style>
  .feature-area {
    height: 100%;
    min-width: 0;
    min-height: 0;
  }

  .feature-area.with-side-navigation {
    display: grid;
    grid-template-columns: 200px minmax(0, 1fr);
  }

  .main-area {
    min-width: 0;
    min-height: 0;
    background: var(--color-background);
  }

  .main-area p {
    color: var(--color-text-muted);
    font-size: 14px;
    margin: 0;
    padding: 8px 16px 16px;
  }
</style>
