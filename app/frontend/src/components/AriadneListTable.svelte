<script lang="ts" generics="T">
  import type { Snippet } from "svelte";

  type Column<T> = {
    key: string;
    header: string;
    cell: Snippet<[T]>;
    width?: string;
    align?: "left" | "center" | "right";
  };

  export let data: T[];
  export let getKey: (item: T) => string | number;
  export let columns: Column<T>[];
</script>

<div class="list-table">
  <div class="list-table-header">
    <span class="list-count">全 {data.length} 件</span>
  </div>

  <div class="table-scroll">
    <table>
      <colgroup>
        {#each columns as column (column.key)}
          <col style:width={column.width} />
        {/each}
      </colgroup>

      <thead>
        <tr>
          {#each columns as column (column.key)}
            <th>{column.header}</th>
          {/each}
        </tr>
      </thead>

      <tbody>
        {#each data as item (getKey(item))}
          <tr>
            {#each columns as column (column.key)}
              <td style:text-align={column.align ?? "left"}>
                {@render column.cell(item)}
              </td>
            {/each}
          </tr>
        {:else}
          <tr class="empty-row">
            <!-- TODO: L02-⑥ 空データ時の表示文言・指定方法 -->
            <td colspan={columns.length}>No tasks.</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<style>
  .list-table {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    min-height: 0;
  }

  /* TODO: L02-⑥ ページタイトル・件数表示・トグルとの配置関係 */
  .list-table-header {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    margin-bottom: 6px;
    min-height: 19px;
  }

  .list-count {
    font-size: 0.75em;
    color: var(--color-text-muted);
    white-space: nowrap;
  }

  .table-scroll {
    flex: 1;
    min-height: 0;
    overflow: auto;
  }

  table {
    width: 100%;
    min-width: 430px; /* TODO: L02-⑥ 最小幅・横スクロールの扱い */
    table-layout: fixed;
    border-collapse: collapse;
    border: 1px solid #d8dde5;
  }

  th,
  td {
    padding: 4px 10px;
    border-bottom: 1px solid #d8dde5;
    font-size: 12px;
    vertical-align: middle;
  }

  th {
    position: sticky; /* TODO: L02-⑦ ヘッダー固定方式 */
    top: 0;
    z-index: 1;
    background: #f4f6f8;
    text-align: center;
  }

  .empty-row td {
    height: 120px;
    text-align: center !important; /* TODO: L02-⑥ 空データ行の表示・スタイル優先順位 */
    color: var(--color-text-muted);
  }

  /* TODO: L02-④ hoverと選択状態の関係 */
  tbody tr:hover {
    outline: 1px solid #4f7fd3;
    outline-offset: -1px;
  }

  /*
  * TODO: L02-⑤ 条件付き行装飾の正式設計時に置き換える。
  * TEMP: Task固有のstatus-badgeを利用してDONE行を判定する。
  */
  tbody tr:has(:global(.status-badge[data-status="DONE"])) {
    background: #e5e7eb;
  }
</style>
