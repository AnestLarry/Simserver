<script lang="ts">
  import { Card, Button } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let showSize: Boolean;
  export let list: Array<any>;
  export let callback: (item: any) => void;
  export let isFile: boolean;
  const dispatch = createEventDispatcher();

  function size(n: number): string {
    let r = [n.toFixed(2), "MB"];
    if (n < 1 || n > 1024) {
      r =
        n < 1 ? [(n * 1024).toFixed(2), " KB"] : [(n / 1024).toFixed(2), "GB"];
    }
    return r.join(" ");
  }
  function download(item: any) {
    dispatch("download", item);
  }
  function copy(item: any) {
    dispatch("copy", item);
  }
</script>

{#each list as LSItem}
  <Card
    class="card"
    on:click={() => {
      if (!isFile) callback(LSItem);
    }}
  >
    <h5 class="font-bold card-title">{LSItem.Name}</h5>
    <p>{new Date(LSItem.ModTime).toLocaleString()}</p>
    {#if showSize}
      <p>{size(LSItem.Size)}</p>
    {/if}
    {#if isFile}
      <div class="button-group">
        <Button on:click={() => callback(LSItem)}>View</Button>
        <Button on:click={() => download(LSItem)}>Download</Button>
        <Button on:click={() => copy(LSItem)}>Copy Link</Button>
      </div>
    {/if}
  </Card>
{/each}

<style>
  .card {
    margin-bottom: 1rem;
  }
  .card-title {
    word-break: break-all;
  }
  .button-group {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
  }
  @media screen and (min-width: 1101px) {
    .card {
      flex: 1 1 calc(31% - 1rem);
      max-width: calc(31% - 1rem);
    }
  }
  @media screen and (max-width: 1100px) and (min-width: 701px) {
    .card {
      flex: 1 1 calc(40% - 1rem);
      max-width: calc(40% - 1rem);
    }
  }
  @media screen and (max-width: 700px) {
    .card {
      flex: 1 1 calc(85% - 1rem);
      min-width: calc(80% - 1rem);
    }
  }
</style>
