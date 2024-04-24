<script>
  export let visible = true;
</script>

{#if visible}
  <slot />
{/if}
