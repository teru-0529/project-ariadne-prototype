<script lang="ts">
  import SideNavigation from "../../components/SideNavigation.svelte";
  import TaskList from "./TaskList.svelte";

  import { TaskScreen, TaskContext } from "./TaskContext";

  // INFO: Side Navigation に渡すリスト
  const navigationItems: { value: TaskScreen; label: string }[] = [
    { value: TaskScreen.List, label: "List" },
    { value: TaskScreen.Create, label: "Create" },
  ];

  // INFO: 選択中の画面
  let currentTaskOperation: TaskScreen | null = null;

  // INFO: タスクコンテキスト
  let context = new TaskContext();
</script>

<div class="feature-area with-side-navigation">
  <SideNavigation
    label="Task"
    items={navigationItems}
    current={currentTaskOperation}
    onSelect={async (value) => {
      // INFO: value が navigationItems に存在すること
      if (!navigationItems.some((item) => item.value === value)) {
        console.error("TaskManagement: Invalid screen", value);
        throw new Error("予期しないエラーが発生しました。");
      }
      const screen = value as TaskScreen;

      await context.prepare(screen);
      context.commit(screen);
      currentTaskOperation = screen;
    }}
  />

  <main class="main-area">
    {#if currentTaskOperation === TaskScreen.List}
      <TaskList {context} />
    {:else if currentTaskOperation === TaskScreen.Create}
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
    overflow: hidden;

    padding: 12px 20px 20px;
    background: var(--color-background);
  }

  .main-area p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 16px;
  }
</style>
