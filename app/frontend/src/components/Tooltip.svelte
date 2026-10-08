<script lang="ts">
  import { tick } from "svelte";

  export let text: string;

  let visible = false;
  let left = 0;
  let top = 0;
  let tooltipElement: HTMLDivElement;

  // Use an action so this display-only element needs no keyboard interaction.
  function hoverTooltip(node: HTMLDivElement) {
    let active = true;

    async function show() {
      if (node.scrollWidth <= node.clientWidth) return;
      visible = true;
      await tick();
      if (!active || !visible || !tooltipElement) return;

      const rect = node.getBoundingClientRect();
      const width = tooltipElement.offsetWidth;
      const height = tooltipElement.offsetHeight;
      left = Math.max(8, Math.min(rect.left, window.innerWidth - width - 8));
      const above = rect.top - height - 8;
      top = above >= 8 ? above : Math.max(8, Math.min(rect.bottom + 8, window.innerHeight - height - 8));
    }

    function hide() {
      visible = false;
    }

    node.addEventListener("pointerenter", show);
    node.addEventListener("pointerleave", hide);
    window.addEventListener("scroll", hide, true);
    window.addEventListener("resize", hide);

    return {
      destroy() {
        active = false;
        node.removeEventListener("pointerenter", show);
        node.removeEventListener("pointerleave", hide);
        window.removeEventListener("scroll", hide, true);
        window.removeEventListener("resize", hide);
      },
    };
  }
</script>

<div class="title" use:hoverTooltip>{text}</div>

{#if visible}
  <div bind:this={tooltipElement} class="tooltip" style:left="{left}px" style:top="{top}px">{text}</div>
{/if}

<style>
  .title {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .tooltip {
    position: fixed;
    z-index: 1000;
    max-width: min(360px, calc(100vw - 16px));
    box-sizing: border-box;
    padding: 8px 12px;
    border: 1px solid #d8dde5;
    border-radius: 6px;
    background: #fff;
    color: #263245;
    box-shadow: 0 3px 12px rgba(25, 40, 65, 0.16);
    font-size: 12px;
    line-height: 1.5;
    overflow-wrap: anywhere;
    pointer-events: none;
  }
</style>
