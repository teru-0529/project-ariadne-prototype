<script lang="ts">
  type Props = {
    type: "error" | "success";
    message: string;
  };

  let { type, message }: Props = $props();
</script>

<div class:error={type === "error"} class:success={type === "success"} class="toast" role="alert">
  <span>{message}</span>
</div>

<style>
  .toast {
    position: fixed;
    top: 38px;
    right: 24px;
    z-index: 1000;

    box-sizing: border-box;
    width: 210px;
    max-width: calc(100vw - 48px);
    padding: 8px;

    font-size: 10px;
    overflow-wrap: anywhere;
    line-height: 1.5;

    border: 1px solid;
    border-radius: 4px;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.25);
  }

  .error {
    background: #b42318;
    color: #ffffff;
    border-color: #912018;
  }

  .success {
    background: #315fc0;
    color: #ffffff;
    border-color: #244b9c;
  }
</style>
