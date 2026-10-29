<script lang="ts">
  import AriadneField from "../../components/AriadneField.svelte";
  import type { FieldVariant } from "../../components/AriadneField.svelte";
  import { TASK_STATUSES } from "./TaskValues";
  import type { Validity } from "../../lib/validity";

  type Props = {
    variant: FieldVariant;
    required?: boolean;
    value: string;
  };

  let { variant, required = false, value }: Props = $props();

  // INFO: 値
  const initialValue = () => value;
  let currentValue = $state(initialValue());

  // INFO: Validity
  let validity = $state<Validity>("UNKNOWN");

  // FUNCTION: 値変更イベント
  function handleValueChange(nextValue: string) {
    currentValue = nextValue;
    validity = "UNKNOWN";
  }

  // TODO: ダミー
  let focusOutCount = $state(0);

  // FUNCTION: フォーカスアウトイベント(グループ内のボタンクリックでもフォーカスアウトと判断される)
  function handleFocusOut(event: FocusEvent) {
    const group = event.currentTarget;

    if (group instanceof HTMLElement && event.relatedTarget instanceof Node && group.contains(event.relatedTarget)) {
      return;
    }

    focusOutCount += 1;
  }

  const label = "ステータス";
</script>

<AriadneField {variant} {label} {required}>
  {#if variant === "edit"}
    <div class="status-group" role="radiogroup" aria-label="ステータス" onfocusout={handleFocusOut}>
      {#each TASK_STATUSES as status}
        <label class="status-option">
          <input
            type="radio"
            name="task-status"
            value={status}
            checked={currentValue === status}
            onchange={() => handleValueChange(status)}
          />
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <span
            class="status-badge"
            class:unselected={currentValue !== status}
            data-status={status}
            onmousedown={(event) => event.preventDefault()}
          >
            {status}
          </span>
        </label>
      {/each}
    </div>
  {:else}
    <span class="status-badge" class:cell={variant === "cell"} data-status={currentValue}>
      {currentValue}
    </span>
  {/if}
  <!-- TODO: ダミー -->
  {#if variant === "edit"}
    <div>validity：{validity}</div>
    <div>Focus Out検知：{focusOutCount}</div>
  {/if}
  <!-- TODO: ダミー -->
</AriadneField>

<style>
  .status-group {
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }

  .status-option {
    position: relative;
    display: inline-flex;
    cursor: pointer;
  }

  /* radioは操作可能なまま、見た目だけ隠す */
  .status-option input {
    position: absolute;
    width: 1px;
    height: 1px;
    opacity: 0;
  }

  /* radioにFocusがあるとき、対応するバッジを強調 */
  .status-option input:focus-visible + .status-badge {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: 1px;
  }

  .status-badge {
    display: block;
    width: fit-content;
    min-width: 52px;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
    line-height: 16px;
    text-align: center;
  }

  .status-badge.cell {
    margin-inline: auto;
  }

  .status-badge[data-status="TODO"] {
    background: #eee4cc;
    color: #80652d;
  }

  .status-badge[data-status="DOING"] {
    background: #dbeafe;
    color: #2563eb;
  }

  .status-badge[data-status="DONE"] {
    background: #d1d5db;
    color: #4b5563;
  }

  /* 非選択時はStatus固有色より優先 */
  .status-badge.unselected {
    background: #ffffff;
    color: #aeb8c5;
    box-shadow: inset 0 0 0 1px #d1d5db;
  }
</style>
