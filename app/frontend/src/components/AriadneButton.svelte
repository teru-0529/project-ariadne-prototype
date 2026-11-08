<script lang="ts">
  type Props = {
    label: string;
    variant?: "primary" | "secondary";
    disabled?: boolean;
    type?: "button" | "submit";
    onclick?: () => void;
  };

  let { label, variant = "primary", disabled = false, type = "button", onclick }: Props = $props();
</script>

<button {type} {disabled} class:primary={variant === "primary"} class:secondary={variant === "secondary"} {onclick}>
  {label}
</button>

<style>
  button {
    min-height: 24px;
    padding: 0 8px;
    border: 1px solid transparent;
    border-radius: 4px;
    font: inherit;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
  }

  button.primary {
    background: #315fba;
    color: #fff;
  }

  button.secondary {
    background: #fff;
    color: #354052;
    border-color: #bfc5ce;
  }

  button:focus-visible {
    outline: 2px solid rgba(59, 111, 189, 0.65);
    outline-offset: 2px;
  }

  button:disabled {
    background: #e7ebf0;
    color: #9aa3b1;
    border-color: transparent;
    cursor: not-allowed;
  }
</style>
