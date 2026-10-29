<script lang="ts">
  import AriadneField from "../../components/AriadneField.svelte";
  import type { FieldVariant } from "../../components/AriadneField.svelte";
  import type { Validity } from "../../lib/validity";

  type Props = {
    variant: FieldVariant;
    required?: boolean;
    value: number;
  };

  let { variant, required = false, value }: Props = $props();

  // INFO: 値
  const initialValue = () => value;
  let currentValue = $state(initialValue());

  // INFO: Validity
  let validity = $state<Validity>("UNKNOWN");

  // FUNCTION: 値変更イベント
  function handleValueChange(nextValue: number) {
    currentValue = nextValue;
    validity = "UNKNOWN";
  }

  // TODO: ダミー
  let focusOutCount = $state(0);

  // FUNCTION: フォーカスアウトイベント
  function handleFocusOut() {
    focusOutCount += 1;
  }

  const label = "優先度";
</script>

<AriadneField {variant} {label} {required}>
  {#if variant === "edit"}
    <div
      class="priority priority-edit"
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

  <!-- TODO: ダミー -->
  {#if variant === "edit"}
    <div>validity：{validity}</div>
    <div>Focus Out検知：{focusOutCount}</div>
  {/if}
  <!-- TODO: ダミー -->
</AriadneField>

<style>
  .priority {
    font-size: 16px;
    line-height: 1;
    transform: translateY(-1px);
    display: inline-block;
    white-space: nowrap;
  }

  .priority-filled {
    color: #dc8a18;
  }

  .priority-empty {
    color: #aeb8c5;
  }

  .priority-edit {
    display: inline-flex;
    align-items: center;
    gap: 1px;
    padding: 4px 6px;
    border-radius: 4px;
    outline: none;
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
    outline-offset: -2px;
  }
</style>
