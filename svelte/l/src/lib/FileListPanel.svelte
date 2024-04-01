<script lang="ts">
  import { client, http } from "../utils";
  import { Card } from "flowbite-svelte";
  export let panel: Panel;
  let flpc: FileListPanelConfig = {
    fileList: [],
    folderList: [],
    urlStack: [],
    listenIndex: -1,
    fresh: (f: boolean) => {
      if (panel.pageMode !== "List") {
        return;
      }
      var urlStack_localStorage = localStorage.getItem("urlStack");
      if (urlStack_localStorage === null)
        throw "fls.fresh: urlStack_localStorage is null";
      var urlStack = JSON.parse(urlStack_localStorage);
      flpc.urlStack = urlStack;
      if (f) {
        var res: LSResponse = JSON.parse(
          http.Get(
            panel.baseUrl +
              "/api/dl/ls/" +
              (urlStack.length > 1 ? urlStack.join("/") : urlStack[0])
          )
        );
        flpc.fileList = res.fileList;
        flpc.folderList = res.folderList;
      } else {
        flpc.fileList.sort(client.sortFunction(panel.sortedBy));
        flpc.folderList.sort(client.sortFunction(panel.sortedBy));
      }
    },
  };
  function init() {
    flpc.fresh(true);
    flpc.listenIndex = panel.workUrlListening.push(flpc.fresh);
  }
  init();
</script>

<div class="file-list-panel">
  <details>
    <summary>
      <h3>Folders:</h3>
    </summary>
    <div class="card-container">
      {#if flpc.folderList && flpc.folderList.length > 0}
        {#each flpc.folderList as LSItem}
          <Card
            class="card"
            on:click={() => {
              panel.pushUrlStack(LSItem.Name);
            }}
          >
            <h5 class="font-bold">{LSItem.Name}</h5>
            <p>{new Date(LSItem.ModTime).toLocaleString()}</p>
          </Card>
        {/each}
      {/if}
    </div>
  </details>
  <div style="clear: both"></div>
  <details>
    <summary>
      <h3>Files:</h3>
    </summary>
    <div class="card-container">
      {#if flpc.fileList && flpc.fileList.length > 0}
        {#each flpc.fileList as LSItem}
          <Card
            class="card"
            href={panel.baseUrl +
              "/api/dl/n/" +
              (flpc.urlStack.length == 1 && flpc.urlStack[0] === "/"
                ? "/" + LSItem.Name
                : flpc.urlStack.join("/") + "/" + LSItem.Name)}
          >
            <h5 class="font-bold">{LSItem.Name}</h5>
            <p>{new Date(LSItem.ModTime).toLocaleString()}</p>
            <p>{LSItem.Size} MB</p>
          </Card>
        {/each}
      {/if}
    </div>
  </details>
</div>

<style>
  .file-list-panel {
    margin: 1em;
  }
  .card-container {
    margin: 1em;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }
  .card {
    margin-bottom: 1rem;
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
