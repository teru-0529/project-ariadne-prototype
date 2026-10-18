<script lang="ts" generics="T">
  import type { Snippet } from "svelte";

  type Column<T> = {
    key: string;
    header: string;
    cell: Snippet<[T]>;
    width?: `${number}px`;
    align?: "left" | "center" | "right";
  };

  export let minWidth: `${number}px` | undefined = undefined;
  export let data: T[];
  export let getKey: (item: T) => string | number;
  export let columns: Column<T>[];
  export let getRowVariant: ((item: T) => "muted" | undefined) | undefined = undefined;
  export let onRowActivate: ((item: T) => void) | undefined = undefined;
</script>

<div class="list-table">
  <div class="list-table-header">
    <span class="list-count">全 {data.length} 件</span>
  </div>

  <div class="table-scroll">
    <table style:min-width={minWidth}>
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
          <tr
            data-variant={getRowVariant?.(item)}
            class:interactive={onRowActivate !== undefined}
            tabindex={onRowActivate ? 0 : undefined}
            onclick={() => onRowActivate?.(item)}
            onkeydown={(event) => {
              if (event.key === "Enter") {
                onRowActivate?.(item);
              }
            }}
          >
            {#each columns as column (column.key)}
              <td style:text-align={column.align ?? "left"}>
                {@render column.cell(item)}
              </td>
            {/each}
          </tr>
        {:else}
          <tr class="empty-row">
            <td colspan={columns.length}>データがありません</td>
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
    height: fit-content;
    max-height: 100%;

    min-width: 0;
    min-height: 0;
  }

  .list-table-header {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    margin-bottom: 0px;
    min-height: 15px;
  }

  .list-count {
    font-size: 0.75em;
    color: var(--color-text-muted);
    white-space: nowrap;
  }

  .table-scroll {
    min-width: 0;
    min-height: 0;
    overflow: auto;
    flex-shrink: 1;
  }

  table {
    width: 100%;
    table-layout: fixed;
    border-collapse: separate;
    border-spacing: 0;
    border-left: 1px solid #d8dde5;
    border-right: 1px solid #d8dde5;
  }

  th,
  td {
    padding: 4px 10px;
    border-bottom: 1px solid #d8dde5;
    font-size: 12px;
    vertical-align: middle;
  }

  th {
    position: sticky;
    top: 0;
    z-index: 1;
    background: #e9edf3;
    text-align: center;

    box-shadow: inset 0 1px 0 #e3e5ea;
  }

  .empty-row td {
    height: 120px;
    text-align: center;
    color: var(--color-text-muted);
  }

  /* INFO: muted：mutedタグのある行を、背景色で示す */
  tbody tr[data-variant="muted"] {
    background: #e5e7eb;
  }

  /* INFO: Hover：操作可能な行だけ、背景色で示す */
  tbody tr.interactive:hover {
    background: #edf4ff;
    cursor: pointer;
  }

  /* INFO:  muted + Hover */
  tbody tr.interactive[data-variant="muted"]:hover {
    background: #dbe5f3;
  }

  /* INFO:  Focus：操作可能な行だけ、青い輪郭で示す */
  tbody tr.interactive:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: -2px;
  }
</style>
