<script lang="ts">
  import { Card } from "flowbite-svelte";

  export let showSize: Boolean;
  export let list: Array<any>;
  export let callback: (item: any) => void;
  function size(n: number): string {
    let r = [n.toFixed(2), "MB"];
    if (n < 1 || n > 1024) {
      r =
        n < 1 ? [(n * 1024).toFixed(2), " KB"] : [(n / 1024).toFixed(2), "GB"];
    }
    return r.join(" ");
  }
</script>

{#each list as LSItem}
  <Card
    class="card"
    on:click={() => {
      callback(LSItem);
    }}
  >
    <h5 class="font-bold card-title">{LSItem.Name}</h5>
    <p>{new Date(LSItem.ModTime).toLocaleString()}</p>
    {#if showSize}
      <p>{size(LSItem.Size)}</p>
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
