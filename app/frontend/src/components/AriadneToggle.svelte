<script lang="ts">
  type ToggleOption = {
    value: string;
    label: string;
  };

  type Props = {
    options: ToggleOption[];
    selected: string;
    onSelect: (value: string) => void;
  };

  let { options, selected, onSelect }: Props = $props();

  function handleSelect(value: string) {
    if (value === selected) return;
    onSelect(value);
  }
</script>

<div class="toggle" role="group">
  {#each options as option (option.value)}
    <button
      type="button"
      class:selected={selected === option.value}
      aria-pressed={selected === option.value}
      onclick={() => handleSelect(option.value)}
    >
      {option.label}
    </button>
  {/each}
</div>

<style>
  .toggle {
    display: inline-flex;
    align-items: stretch;
  }

  button {
    width: 64px;
    min-height: 24px;
    padding: 3px 4px;
    border: 1px solid #c9d1dc;
    background: #ffffff;
    color: #66758a;
    font: inherit;
    font-size: 10px;
    font-weight: 400;
    cursor: pointer;
  }

  button.selected {
    position: relative;
    border: 1px solid #3b6fbd;
    background: #edf4ff;
    color: #3267b1;
    font-weight: 600;
  }

  button + button {
    margin-left: -1px;
  }

  button:first-child {
    border-radius: 5px 0 0 5px;
  }

  button:last-child {
    border-radius: 0 5px 5px 0;
  }

  button:focus-visible {
    position: relative;
    z-index: 1;
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: 1px;
  }
</style>
