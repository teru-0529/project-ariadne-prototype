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
    onSuccess: (message: string) => void;
  };

  let { onError, onSuccess }: Props = $props();

  // INFO: 各Pageへの参照
  let taskListPage = $state<TaskListPage>();
  let taskDetailPage = $state<TaskDetailPage>();
  let taskCreatePage = $state<TaskCreatePage>();

  // INFO: Side Navigation に渡すリスト
  const navigationItems: { value: TaskScreen; label: string }[] = [
    { value: TaskScreen.List, label: "List" },
    { value: TaskScreen.Create, label: "Create" },
  ];

  // INFO: 画面再生成キー
  let navigationKey = $state(0);

  // INFO: 画面遷移中
  let isNavigating = false;

  // INFO: 選択中の画面
  let currentPage = $state<TaskScreen | null>(null);

  // INFO: タスクコンテキスト
  let context = new TaskContext();

  // FUNCTION: Side Navigationからの遷移要求
  function handleSideNavSelect(value: string): void {
    // INFO: Navigationに指定した文字列であることをチェック
    if (!navigationItems.some((item) => item.value === value)) {
      handleError(createProcessedError(`TaskManagement: Invalid screen: ${value}`));
      return;
    }

    const screen = value as typeof TaskScreen.List | typeof TaskScreen.Create;
    const request: TaskNavigation = { screen };

    // INFO: Pageに離脱依頼
    leave(() => {
      void navigate(request);
    });
  }

  // FUNCTION: 現在のPageへの離脱問い合わせ
  export function leave(callback: () => void): void {
    if (currentPage === TaskScreen.List) {
      taskListPage!.leave(callback);
    } else if (currentPage === TaskScreen.Detail) {
      taskDetailPage!.leave(callback);
    } else if (currentPage === TaskScreen.Create) {
      taskCreatePage!.leave(callback);
    } else {
      // INFO: 初期状態はPageが存在しない
      callback();
    }
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
      context.commit();
      currentPage = request.screen;
      navigationKey += 1;
    } catch (error) {
      handleError(error);
    } finally {
      isNavigating = false;
    }
  }

  // FUNCTION: 再遷移Callbackの生成
  function createReloadCallback(): () => void {
    const navigation = context.currentNavigation();

    return () => {
      void navigate(navigation);
    };
  }

  // FUNCTION: エラーハンドリング
  function handleError(error: unknown): void {
    onError(toErrorMessage(error));
  }

  // FUNCTION: Task一覧からの遷移要求
  function handleTaskSelect(id: string): void {
    const request: TaskNavigation = { screen: TaskScreen.Detail, id };

    // INFO: Pageに離脱依頼
    leave(() => {
      void navigate(request);
    });
  }

  // FUNCTION: Task登録成功後の再遷移
  function handleTaskCreated(): void {
    void navigate({ screen: TaskScreen.Create });
  }
</script>

<div class="feature-area with-side-navigation">
  <SideNavigation label="Task" items={navigationItems} current={currentPage} onSelect={handleSideNavSelect} />

  <main class="main-area">
    {#key navigationKey}
      {#if currentPage === TaskScreen.List}
        <TaskListPage bind:this={taskListPage} {context} onTaskSelect={handleTaskSelect} />
      {:else if currentPage === TaskScreen.Detail}
        <TaskDetailPage bind:this={taskDetailPage} {context} />
      {:else if currentPage === TaskScreen.Create}
        <TaskCreatePage bind:this={taskCreatePage} {onSuccess} {onError} onCreated={handleTaskCreated} />
      {:else}
        <p>Select an operation.</p>
      {/if}
    {/key}
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
