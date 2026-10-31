<script lang="ts">
  import AriadneField from "../../components/AriadneField.svelte";
  import type { FieldVariant } from "../../components/AriadneField.svelte";
  import { TASK_STATUSES } from "./TaskValues";

  import type { Validation } from "../../lib/validity";

  type Props = {
    variant: FieldVariant;
    required?: boolean;
    value: string;
  };

  let { variant, required = false, value }: Props = $props();

  // INFO: 値
  const initialValue = () => value;
  let currentValue = $state(initialValue());
  let validation = $state<Validation>({ validity: "UNKNOWN" });

  // INFO: エラー表示の活性状態
  let errorActivated = $state(false);

  // INFO:TODO: 表示用エラーメッセージ
  // let displayError = "テスト用ダミーメッセージ";
  let displayError = $derived(errorActivated && validation.validity === "INVALID" ? validation.errorMessage : null);

  // FUNCTION: 値変更イベント
  function handleValueChange(nextValue: string) {
    currentValue = nextValue;
    validation = { validity: "UNKNOWN" };
  }

  // TODO: ダミー
  let focusOutCount = $state(0);

  // FUNCTION: フォーカスアウトイベント(グループ内のボタンクリックでもフォーカスアウトと判断される)
  function handleFocusOut(event: FocusEvent) {
    const group = event.currentTarget;

    if (group instanceof HTMLElement && event.relatedTarget instanceof Node && group.contains(event.relatedTarget)) {
      return;
    }

    errorActivated = true;
    focusOutCount += 1;
  }

  const label = "ステータス";
</script>

<AriadneField {variant} {label} {required} error={displayError}>
  {#if variant === "edit"}
    <div
      class="status-group"
      class:invalid={!!displayError}
      role="radiogroup"
      aria-label="ステータス"
      onfocusout={handleFocusOut}
    >
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
</AriadneField>
<!-- TODO: ダミー -->
{#if variant === "edit"}
  <div>
    validity：{validation.validity}
    {#if validation.validity === "INVALID"}
      / errorMessage：{validation.errorMessage}
    {/if}
  </div>
  <div>errorActivated：{errorActivated}</div>
  <div>Focus Out検知：{focusOutCount}</div>
{/if}

<!-- TODO: ダミー -->

<style>
  .status-group {
    display: inline-flex;
    align-items: center;
    margin-top: 5px;
    margin-bottom: 3px;
    margin-left: 4px;
    gap: 6px;
    transform: translateY(-1px);
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
    outline-offset: 2px;
  }

  .status-group.invalid {
    outline: 2px solid #c62828;
    outline-offset: 4px;
    border-radius: 1px;
  }

  .status-badge {
    display: block;
    width: fit-content;
    min-width: 52px;
    padding: 1px 6px;
    border-radius: 3px;
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
