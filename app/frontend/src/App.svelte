<script lang="ts">
  import { onMount } from "svelte";
  import { AppService } from "../bindings/ariadne/backend/services";

  let appName = $state("");

  onMount(async () => {
    try {
      appName = await AppService.GetAppName();
    } catch (error) {
      console.error("Error initializing AppService:", error);
    }
  });
</script>

<main class="container">
  <h1 class="title">{appName}</h1>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: sans-serif;
  }

  main {
    padding: 2rem;
    color: white;
  }
</style>
