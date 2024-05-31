<script lang="ts">
  import FileCardList from "../lib/FileCardList.svelte";
  import Toggleable from "../lib/Toggleable.svelte";
  import { client, http } from "../utils";
  import { Input } from "flowbite-svelte";
  export let panel: Panel;
  let flpc: FileListPanelConfig = {
    fileList: [],
    folderList: [],
    urlStack: [],
    listenIndex: -1,
    filterCond: "",
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
        flpc.fileList = [];
        flpc.folderList = [];
        var res: LSResponse = JSON.parse(
          http.Get(
            panel.baseUrl +
              "/api/dl/ls/" +
              (urlStack.length > 1 ? urlStack.join("/") : urlStack[0])
          )
        );
        if (res.fileList !== null)
          flpc.fileList = res.fileList.sort(
            client.sortFunction(panel.sortedBy)
          );
        if (res.folderList !== null)
          flpc.folderList = res.folderList.sort(
            client.sortFunction(panel.sortedBy)
          );
      } else {
        flpc.fileList.sort(client.sortFunction(panel.sortedBy));
        flpc.folderList.sort(client.sortFunction(panel.sortedBy));
      }
    },
    fileListFiltered: function (): LSItem[] {
      return callbacks.filterCallback(flpc.fileList);
    },
    folderListFiltered: function (): LSItem[] {
      return callbacks.filterCallback(flpc.folderList);
    },
  };
  function init() {
    flpc.fresh(true);
    flpc.listenIndex = panel.workUrlListening.push(flpc.fresh);
  }
  init();
  const callbacks = {
    folderCallback: (item: LSItem) => {
      panel.pushUrlStack(item.Name);
    },
    fileCallback: (item: LSItem) => {
      const url =
        panel.baseUrl +
        "/api/dl/n/" +
        (flpc.urlStack.length == 1 && flpc.urlStack[0] === "/"
          ? "/" + item.Name
          : flpc.urlStack.join("/") + "/" + item.Name);
      window.open(url, "_blank");
    },
    filterCallback: (list: LSItem[]): LSItem[] => {
      return list.filter((x) => x.Name.indexOf(flpc.filterCond) > -1);
    },
  };
</script>

<div class="file-list-panel">
  <Input bind:value={flpc.filterCond} placeholder="filter" />
  <details>
    <summary>
      <h3>Folders:</h3>
    </summary>
    <div class="card-container">
      <Toggleable visible={flpc.folderList && flpc.folderList.length > 0}>
        <FileCardList
          list={flpc.folderListFiltered()}
          callback={callbacks.folderCallback}
          showSize={false}
        />
      </Toggleable>
    </div>
  </details>
  <div style="clear: both"></div>
  <details>
    <summary>
      <h3>Files:</h3>
    </summary>
    <div class="card-container">
      <Toggleable visible={flpc.fileList && flpc.fileList.length > 0}>
        <FileCardList
          list={flpc.fileListFiltered()}
          showSize={true}
          callback={callbacks.fileCallback}
        />
      </Toggleable>
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
</style>
