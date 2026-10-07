<script lang="ts">
  type NavigationItem = {
    value: string;
    label: string;
  };

  export let label: string;
  export let items: NavigationItem[];
</script>

<aside class="side-navigation">
  <nav aria-label={label}>
    {#each items as item}
      <button type="button">
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
</style>
