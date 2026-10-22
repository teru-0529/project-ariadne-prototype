<script lang="ts">
  import SideNavigation from "../../components/AriadneSideNavigation.svelte";
  import TaskListPage from "./TaskListPage.svelte";
  import TaskDetailPage from "./TaskDetailPage.svelte";
  import TaskCreatePage from "./TaskCreatePage.svelte";

  import { TaskScreen, TaskContext, type TaskNavigation } from "./TaskContext";
  import { toErrorMessage } from "../../lib/errors/ProcessedError";
  import { createProcessedError } from "../../lib/errors/ProcessedError";

  type Props = {
    onError: (message: string) => void;
  };

  let { onError }: Props = $props();

  // INFO: Side Navigation に渡すリスト
  const navigationItems: { value: TaskScreen; label: string }[] = [
    { value: TaskScreen.List, label: "List" },
    { value: TaskScreen.Create, label: "Create" },
  ];

  // INFO: 画面遷移中
  let isNavigating = false;

  // INFO: 選択中の画面
  let currentTaskOperation = $state<TaskScreen | null>(null);

  // INFO: タスクコンテキスト
  let context = new TaskContext();

  // FUNCTION: Side Navigationからの選択
  function handleSideNavSelect(value: string): void {
    // INFO: Navigationに指定した文字列であることをチェック
    if (!navigationItems.some((item) => item.value === value)) {
      handleError(createProcessedError(`TaskManagement: Invalid screen: ${value}`));
      return;
    }

    const screen = value as typeof TaskScreen.List | typeof TaskScreen.Create;
    void navigate({ screen });
  }

  // FUNCTION: Task一覧からの選択
  function handleTaskSelect(id: string): void {
    void navigate({
      screen: TaskScreen.Reference,
      id,
    });
  }

  // FUNCTION: 共通の画面遷移
  async function navigate(request: TaskNavigation): Promise<void> {
    // INFO: 二重操作防止
    if (isNavigating) {
      return;
    }

    // INFO: 遷移処理
    isNavigating = true;
    try {
      await context.prepare(request);
      context.commit(request.screen);
      currentTaskOperation = request.screen;
    } catch (error) {
      handleError(error);
    } finally {
      isNavigating = false;
    }
  }

  // FUNCTION: エラーハンドリング
  function handleError(error: unknown): void {
    onError(toErrorMessage(error));
  }
</script>

<div class="feature-area with-side-navigation">
  <SideNavigation label="Task" items={navigationItems} current={currentTaskOperation} onSelect={handleSideNavSelect} />

  <main class="main-area">
    {#if currentTaskOperation === TaskScreen.List}
      <TaskListPage {context} onTaskSelect={handleTaskSelect} />
    {:else if currentTaskOperation === TaskScreen.Reference}
      <TaskDetailPage {context} />
    {:else if currentTaskOperation === TaskScreen.Create}
      <TaskCreatePage />
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
