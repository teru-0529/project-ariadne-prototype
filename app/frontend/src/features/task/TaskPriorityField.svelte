<script lang="ts">
  import AriadneField from "../../components/AriadneField.svelte";
  import type { FieldVariant } from "../../components/AriadneField.svelte";

  import type { Validation } from "../../lib/validity";

  type Props = {
    variant: FieldVariant;
    required?: boolean;
    value: number;
  };

  let { variant, required = false, value }: Props = $props();

  // INFO: 値
  const initialValue = () => value;
  let currentValue = $state(initialValue());
  let validation = $state<Validation>({ validity: "UNKNOWN" });

  // INFO: エラー表示の活性状態
  let errorActivated = $state(false);

  // INFO: 表示用エラーメッセージ
  let displayError = $derived(errorActivated && validation.validity === "INVALID" ? validation.errorMessage : null);

  // FUNCTION: Validation評価
  function validate() {
    validation = { validity: "VALID" };
  }

  // FUNCTION: 値変更イベント
  function handleValueChange(nextValue: number) {
    currentValue = nextValue;
    validation = { validity: "UNKNOWN" };
  }

  // TODO: ダミー
  let focusOutCount = $state(0);

  // FUNCTION: フォーカスアウトイベント
  function handleFocusOut() {
    errorActivated = true;
    validate();
    focusOutCount += 1;
  }

  const label = "優先度";
</script>

<!-- TODO: ダミー -->
{#if variant === "edit"}
  <div class="callback-debug">
    validity：{validation.validity}
    {#if validation.validity === "INVALID"}
      / errorMessage：{validation.errorMessage}
    {/if}
    | errorActivated：{errorActivated}
    | Focus Out検知：{focusOutCount}
  </div>
{/if}

<!-- TODO: ダミー -->
<AriadneField {variant} {label} {required} error={displayError}>
  {#if variant === "edit"}
    <div
      class="priority priority-edit"
      class:invalid={!!displayError}
      role="spinbutton"
      tabindex="0"
      aria-label="優先度"
      aria-valuemin={1}
      aria-valuemax={5}
      aria-valuenow={currentValue}
      onkeydown={(event) => {
        if (event.key === "ArrowRight" || event.key === "ArrowUp") {
          event.preventDefault();
          handleValueChange(Math.min(5, currentValue + 1));
        } else if (event.key === "ArrowLeft" || event.key === "ArrowDown") {
          event.preventDefault();
          handleValueChange(Math.max(1, currentValue - 1));
        }
      }}
      onfocusout={handleFocusOut}
    >
      {#each [1, 2, 3, 4, 5] as star}
        <span
          class="priority-star"
          class:priority-filled={star <= currentValue}
          class:priority-empty={star > currentValue}
          onclick={() => handleValueChange(star)}
          aria-hidden="true"
        >
          {star <= currentValue ? "★" : "☆"}
        </span>
      {/each}
    </div>
  {:else}
    <span class="priority" role="img" aria-label={`優先度 ${value}`}>
      <span class="priority-filled">{"★".repeat(value)}</span><span class="priority-empty">{"☆".repeat(5 - value)}</span
      >
    </span>
  {/if}
</AriadneField>

<style>
  .priority {
    font-size: 16px;
    line-height: 1;
    transform: translateY(-1px);
    display: inline-block;
    white-space: nowrap;
  }

  .priority-edit {
    display: inline-flex;
    align-items: center;
    transform: none;
    gap: 1px;
    padding: 4px 6px;
    border-radius: 4px;
    outline: none;
  }

  .priority-filled {
    color: #dc8a18;
  }

  .priority-empty {
    color: #aeb8c5;
  }

  .priority-star {
    cursor: pointer;
    font-size: 17px;
    line-height: 1;
    transform: translateY(-1px);
    user-select: none;
  }

  .priority-edit:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: 0px;
  }

  /* INFO: エラー時：Focusの有無にかかわらず赤枠 */
  .priority-edit.invalid,
  .priority-edit.invalid:focus-visible {
    outline: 2px solid #c62828;
    outline-offset: 0px;
  }

  /* TODO: debug */
  .callback-debug {
    padding: 8px 6px;
    font-size: 10px;
    line-height: 1.8;
  }
</style>
