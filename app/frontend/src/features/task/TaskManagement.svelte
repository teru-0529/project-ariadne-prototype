<script lang="ts">
  import SideNavigation from "../../components/SideNavigation.svelte";
  import TaskList from "./TaskList.svelte";

  const navigationItems = [
    { value: "list", label: "List" },
    { value: "create", label: "Create" },
  ];

  let currentTaskOperation: string | null = null;
</script>

<div class="feature-area with-side-navigation">
  <SideNavigation
    label="Task"
    items={navigationItems}
    current={currentTaskOperation}
    onSelect={(value) => {
      currentTaskOperation = value;
    }}
  />

  <main class="main-area">
    {#if currentTaskOperation === "list"}
      <TaskList />
      <!-- {:else if currentTaskOperation === "create"}
      <p>Create Task.</p> -->
    {:else}
      <p>Select an operation.</p>
    {/if}
  </main>
</div>

<style>
  .feature-area {
    height: 100%;
    min-width: 0;
    min-height: 0;
  }

  .feature-area.with-side-navigation {
    display: grid;
    grid-template-columns: 200px minmax(0, 1fr);
  }

  .main-area {
    min-width: 0;
    min-height: 0;
    overflow: hidden;

    padding: 12px 20px 20px;
    background: var(--color-background);
  }

  .main-area p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 16px;
  }
</style>
