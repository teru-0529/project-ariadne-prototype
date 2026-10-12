<script lang="ts">
  import TaskManagement from "./features/task/TaskManagement.svelte";
  import Toast from "./components/Toast.svelte";

  // INFO: Primary Navigation のリスト
  type PrimaryPage = "task" | "config";

  // INFO: 選択中のPrimaryページ
  let currentPage: PrimaryPage = "task";

  // INFO: Toast
  type ToastState = {
    type: "error" | "success";
    message: string;
  };

  // INFO: Toast表示状態
  let toast: ToastState | null = null;

  // INFO: Toastのタイマー（Success用）
  let successTimer: ReturnType<typeof setTimeout> | null = null;

  // FUNCTION: エラーToast
  function notifyError(message: string): void {
    if (successTimer !== null) {
      clearTimeout(successTimer);
      successTimer = null;
    }

    toast = {
      type: "error",
      message,
    };
  }

  // FUNCTION: 完了Toast
  function notifySuccess(message: string): void {
    if (successTimer !== null) {
      clearTimeout(successTimer);
    }

    toast = {
      type: "success",
      message,
    };

    successTimer = setTimeout(() => {
      toast = null;
      successTimer = null;
    }, 3000);
  }

  // FUNCTION: Toast削除（Toast以外の部分をクリックしたときに削除）
  function handleWindowClick(event: MouseEvent): void {
    const target = event.target;

    if (target instanceof Element && target.closest(".toast")) {
      return;
    }

    if (successTimer !== null) {
      clearTimeout(successTimer);
      successTimer = null;
    }

    toast = null;
  }
</script>

<div class="application">
  <header class="application-header">
    <div class="brand">
      <div class="brand-accent" aria-hidden="true">
        <div class="brand-bar"></div>
      </div>

      <div class="brand-title">
        <h1>PROJECT ARIADNE</h1>
      </div>
    </div>

    <nav class="primary-navigation" aria-label="Primary">
      <button
        type="button"
        aria-current={currentPage === "task" ? "page" : undefined}
        on:click={() => (currentPage = "task")}
      >
        Task
      </button>
      <button
        type="button"
        aria-current={currentPage === "config" ? "page" : undefined}
        on:click={() => (currentPage = "config")}
      >
        Config
      </button>
    </nav>
  </header>

  <div class="workspace">
    {#if currentPage === "task"}
      <TaskManagement onError={notifyError} />
    {:else}
      <main class="main-area">
        <p>Configuration is not available in this prototype.</p>
      </main>
    {/if}
  </div>

  <footer class="application-footer">
    <p>&copy; 2026 PROJECT ARIADNE. All rights reserved.</p>
  </footer>
</div>

<!-- INFO: Toast -->
{#if toast !== null}
  <Toast type={toast.type} message={toast.message} />
{/if}

<!-- INFO: Windowクリック時の動き -->
<svelte:window on:click={handleWindowClick} />

<style>
  :global(*) {
    box-sizing: border-box;
  }

  :global(html),
  :global(body),
  :global(#app) {
    width: 100%;
    height: 100%;
    margin: 0;
  }

  :global(body) {
    font-family: "Segoe UI", "Note Sans JP", sans-serif;
  }

  button {
    font: inherit;
  }

  .application {
    --color-background: #ffffff;
    --color-surface: #f7f8fa;
    --color-border: #d9dde5;
    --color-text: #252a33;
    --color-text-muted: #667085;
    --color-accent: #315fba;
    --color-footer: #263b5e;

    display: grid;
    grid-template-rows: 64px minmax(0, 1fr) 32px;

    width: 100%;
    height: 100%;

    color: var(--color-text);
    background: var(--color-background);
  }

  .application-header {
    min-width: 0;
    overflow-x: clip;

    display: flex;
    align-items: stretch;
    justify-content: space-between;

    padding: 0 24px 0 0;

    background: #ffffff;
    border-bottom: 1px solid var(--color-border);
  }

  .brand {
    display: flex;
    align-items: stretch;

    height: 100%;
  }

  /* 左の薄い水色ブロック */
  .brand-accent {
    position: relative;

    width: 48px;
    margin: 0;

    background: #dce8fb;
  }

  /* 濃いブルーの縦棒 */
  .brand-bar {
    position: absolute;

    top: 0;
    right: 0;
    bottom: 0;

    width: 15px;

    background: var(--color-accent);
  }

  .brand-title {
    position: relative;

    display: flex;
    align-items: center;

    margin-left: 12px;
    padding-right: 180px;
  }

  .brand-title h1 {
    margin: 0;

    font-size: 24px;
    font-weight: 600;
    letter-spacing: 0.02em;
    white-space: nowrap;
  }

  /* タイトル下の右へ伸びるライン */
  .brand-bar::after {
    position: absolute;
    bottom: 4px;
    left: 100%;

    width: 360px;
    height: 5px;

    content: "";
    background: linear-gradient(to right, var(--color-accent) 0%, var(--color-accent) 55%, transparent 100%);
  }

  .primary-navigation {
    display: flex;
    align-items: stretch;
    gap: 8px;
  }

  .primary-navigation button {
    position: relative;

    min-width: 72px;
    padding: 0 16px;

    color: var(--color-text-muted);
    background: transparent;
    border: 0;

    cursor: pointer;
  }

  .primary-navigation button[aria-current="page"] {
    color: var(--color-accent);
    font-weight: 600;
  }

  .primary-navigation button[aria-current="page"]::after {
    position: absolute;
    right: 12px;
    bottom: 0;
    left: 12px;

    height: 3px;

    content: "";
    background: var(--color-accent);
  }

  .workspace {
    min-width: 0;
    min-height: 0;
  }

  .main-area {
    padding: 32px;
    background: var(--color-background);
  }

  .main-area p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 16px;
  }

  .application-footer {
    display: flex;
    align-items: center;
    justify-content: center;

    padding: 0 16px;

    color: #eef2f7;
    background: var(--color-footer);

    font-size: 12px;
  }
</style>
