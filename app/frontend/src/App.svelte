<script lang="ts">
  import { onMount } from "svelte";
  import { WML } from "@wailsio/runtime";

  const wailsVersion = "v3.0.0-beta.12";

  onMount(() => {
    // Wire up data-wml-openURL links (logos + footer "Docs" link).
    WML.Reload();
  });
</script>

<main class="container">
  <header class="brand">
    <span class="brand-mark" data-wml-openURL="https://v3.wails.io">
      <img src="/wails.png" class="brand-logo" alt="Wails logo" />
    </span>
    <span class="brand-badge" data-wml-openURL="https://svelte.dev">
      <img src="/svelte.svg" alt="Svelte logo" />
    </span>
  </header>

  <h1 class="title">
    <span class="title-accent">Wails +</span>
    <span class="title-name">
      <span class="title-name-text">Svelte</span>
    </span>
  </h1>
  <p class="subtitle">
    Build beautiful cross-platform apps with Go and Svelte.
  </p>

  <div class="greet">
    <div class="input-box">
      <svg
        class="input-icon"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        aria-hidden="true"
        ><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" /><circle
          cx="12"
          cy="7"
          r="4"
        /></svg
      >
    </div>
  </div>
</main>

<hr class="footer-divider" />
<footer class="footer">
  <span class="footer-version">{wailsVersion}</span>
  <span class="footer-time">
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
      ><circle cx="12" cy="12" r="10" /><polyline
        points="12 6 12 12 16 14"
      /></svg
    >
  </span>
  <a
    class="footer-docs"
    data-wml-openURL="https://v3.wails.io"
    aria-label="Wails documentation"
    >Docs
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
      ><line x1="7" y1="17" x2="17" y2="7" /><polyline
        points="7 7 17 7 17 17"
      /></svg
    >
  </a>
</footer>

<style>
  /* Put your standard CSS here */
</style>
