<script lang="ts">
  import { Card } from "flowbite-svelte";

  export let showSize: Boolean;
  export let list: Array<any>;
  export let callback: (item: any) => void;
</script>

{#each list as LSItem}
  <Card
    class="card"
    on:click={() => {
      callback(LSItem);
    }}
  >
    <h5 class="font-bold">{LSItem.Name}</h5>
    <p>{new Date(LSItem.ModTime).toLocaleString()}</p>
    {#if showSize}
      <p>{LSItem.Size} MB</p>
    {/if}
  </Card>
{/each}
