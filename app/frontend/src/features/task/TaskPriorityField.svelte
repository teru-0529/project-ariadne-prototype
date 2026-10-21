<script lang="ts">
  import AriadneField from "../../components/AriadneField.svelte";
  import type { FieldVariant } from "../../components/AriadneField.svelte";

  type Props = {
    variant: FieldVariant;
    label?: string;
    value: number;
  };

  let { variant, label = "", value }: Props = $props();

  const initialValue = () => value;
  let currentValue = $state(initialValue());
</script>

<AriadneField {variant} {label}>
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
          currentValue = Math.min(5, currentValue + 1);
        } else if (event.key === "ArrowLeft" || event.key === "ArrowDown") {
          event.preventDefault();
          currentValue = Math.max(1, currentValue - 1);
        }
      }}
    >
      {#each [1, 2, 3, 4, 5] as star}
        <span
          class="priority-star"
          class:priority-filled={star <= currentValue}
          class:priority-empty={star > currentValue}
          onclick={() => (currentValue = star)}
          aria-hidden="true"
        >
          {star <= currentValue ? "★" : "☆"}
        </span>
      {/each}
    </div>
  {:else}
    <span class="priority" role="img" aria-label={`優先度 ${currentValue}`}>
      <span class="priority-filled">{"★".repeat(currentValue)}</span><span class="priority-empty"
        >{"☆".repeat(5 - currentValue)}</span
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
    outline: 2px solid #3264bd;
    outline-offset: 2px;
  }
</style>
