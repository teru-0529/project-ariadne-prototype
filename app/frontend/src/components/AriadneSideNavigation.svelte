<script lang="ts">
  type NavigationItem = {
    value: string;
    label: string;
  };

  export let label: string;
  export let items: NavigationItem[];
  export let current: string | null;
  export let onSelect: (value: string) => void;
</script>

<aside class="side-navigation">
  <nav aria-label={label}>
    {#each items as item}
      <button type="button" aria-pressed={current === item.value} on:click={() => onSelect(item.value)}>
        {item.label}
      </button>
    {/each}
  </nav>
</aside>

<style>
  .side-navigation {
    padding: 12px 8px;

    background: var(--color-surface);
    border-right: 1px solid var(--color-border);
  }

  .side-navigation nav {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .side-navigation button {
    width: 100%;
    padding: 7px 12px;

    color: var(--color-text);
    text-align: left;

    background: transparent;
    border: 0;
    border-radius: 3px;

    cursor: pointer;
  }

  .side-navigation button:hover {
    background: #eceff4;
  }

  .side-navigation button[aria-pressed="true"] {
    color: var(--color-primary);
    background: #e8eef9;
    font-weight: 600;
  }

  .side-navigation button[aria-pressed="true"]:hover {
    background: #dce6f7;
  }

  /* INFO: キーボードFocus */
  .side-navigation button:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: -2px;
  }
</style>
