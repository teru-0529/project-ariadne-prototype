<script lang="ts" module>
  // INFO: TYPE: 表示形態
  export type FieldVariant = "reference" | "edit" | "cell";
</script>

<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    variant: FieldVariant;
    label?: string;
    required?: boolean;
    error?: string | null;
    children: Snippet;
  };

  let { variant, label = "", required = false, error = null, children }: Props = $props();
</script>

{#if variant === "cell"}
  {@render children()}
{:else}
  <div class="ariadne-field">
    <div class="field-label">
      {label}
      {#if variant === "edit" && required}
        <span class="required">必須</span>
      {/if}
    </div>

    <div class="field-value" class:editable={variant === "edit"} class:invalid={variant === "edit" && !!error}>
      {@render children()}
    </div>

    <div class="field-error">
      {#if variant === "edit" && error}
        {error}
      {/if}
    </div>
  </div>
{/if}

<style>
  .ariadne-field {
    display: flex;
    flex-direction: column;
    gap: 1px;
    margin-bottom: 8px;
  }

  .field-label {
    font-size: 9px;
  }

  .field-value {
    min-height: 28px;
    padding: 5px 8px;
    border-radius: 4px;
    background: #f3f4f6;
    box-sizing: border-box;
  }

  .field-value.editable {
    background: #fff;
    border: 1px solid #bfc5ce;
  }

  .field-value.editable.invalid {
    border-color: #d94855;
  }

  /* FocusとErrorが重なった場合はErrorを優先 */
  .field-value.editable.invalid:focus-within {
    border-color: #d94855;
    outline-color: rgba(217, 72, 85, 0.2);
  }

  .field-error {
    min-height: 12px;
    padding-left: 5px;
    font-size: 10px;
    color: #c62828;
  }

  .field-label {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 9px;
  }

  .required {
    padding: 2px 5px;
    border-radius: 3px;
    background: #e5f1ea;
    color: #48765d;
    font-size: 9px;
    font-weight: 600;
    line-height: 1.2;
  }

  /* Focus */
  .field-value.editable:focus-within {
    border-color: #3b6fbd;
    outline: 2px solid rgba(59, 111, 189, 0.2);
  }
</style>
