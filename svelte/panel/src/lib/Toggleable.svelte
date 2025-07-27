<script>
  let { visible, children } = $props();
</script>

{#if visible}
  {@render children?.()}
{/if}
