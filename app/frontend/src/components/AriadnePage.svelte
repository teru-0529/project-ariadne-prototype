<script lang="ts">
  import type { Snippet } from "svelte";

  // INFO: ページタイトル
  export let title: string;

  // INFO: ページヘッダー右端の任意コンテンツ
  export let headerAction: Snippet | undefined = undefined;

  // INFO: ページ本文
  export let children: Snippet;
</script>

<div class="ariadne-page">
  <div class="page-header">
    <h2>{title}</h2>

    {#if headerAction}
      <div class="header-action">
        {@render headerAction()}
      </div>
    {/if}
  </div>

  <div class="page-content">
    {@render children()}
  </div>
</div>

<style>
  .ariadne-page {
    display: grid;
    grid-template-rows: auto minmax(0, 1fr);
    width: 100%;
    height: 100%;
    min-height: 0;
    padding: 8px 16px 16px;
    overflow: hidden;
  }

  .page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0;
  }

  .page-header h2 {
    margin: 0;
    font-size: 1em;
    font-weight: 600;
  }

  .page-content {
    min-height: 0;
    min-width: 0;
  }
</style>
