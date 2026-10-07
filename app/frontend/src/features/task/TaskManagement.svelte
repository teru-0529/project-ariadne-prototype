<script lang="ts">
  import SideNavigation from "../../components/SideNavigation.svelte";

  const navigationItems = [
    { value: "list", label: "List" },
    { value: "create", label: "Create" },
  ];
</script>

<div class="feature-area with-side-navigation">
  <SideNavigation label="Task" items={navigationItems} />

  <main class="main-area">
    <p>Select an operation.</p>
  </main>
</div>

<style>
  .feature-area {
    height: 100%;
    min-height: 0;
  }

  .feature-area.with-side-navigation {
    display: grid;
    grid-template-columns: 200px minmax(0, 1fr);
  }

  .main-area {
    padding: 32px;
    background: var(--color-background);
  }

  .main-area p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 16px;
  }
</style>
