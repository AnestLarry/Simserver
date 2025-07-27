<script lang="ts">
  import FileCardList from "../lib/FileCardList.svelte";
  import Toggleable from "../lib/Toggleable.svelte";
  import { client, http } from "../utils";
  import Input from "../lib/components/ui/input.svelte";
  import Button from "../lib/components/ui/button.svelte";
  import Label from "../lib/components/ui/label.svelte";
  import Select from "../lib/components/ui/select.svelte";
  
  let { panel }: { panel: Panel } = $props();
  
  let flpc: FileListPanelConfig = $state({
    fileList: [],
    folderList: [],
    urlStack: [],
    listenIndex: -1,
    filterCond: "",
    showFolders: true,
    showFiles: true,
    columnCount: "4",
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
        try {
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
        } catch (error) {
          console.warn("FileListPanel API call failed (ls feature may be disabled):", error);
          flpc.fileList = [];
          flpc.folderList = [];
        }
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
  });
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
    fileDownloadCallback: (item: LSItem) => {
      const url = 
        panel.baseUrl +
        "/api/dl/n/" +
        (flpc.urlStack.length == 1 && flpc.urlStack[0] === "/"
          ? "/" + item.Name
          : flpc.urlStack.join("/") + "/" + item.Name);
      const a = document.createElement('a');
      a.href = url;
      a.download = item.Name;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
    },
    fileCopyCallback: (item: LSItem) => {
      const url = 
        panel.baseUrl +
        "/api/dl/n/" +
        (flpc.urlStack.length == 1 && flpc.urlStack[0] === "/"
          ? "/" + item.Name
          : flpc.urlStack.join("/") + "/" + item.Name);
      navigator.clipboard.writeText(url);
    },
    filterCallback: (list: LSItem[]): LSItem[] => {
      return list.filter((x) => x.Name.indexOf(flpc.filterCond) > -1);
    },
  };
</script>

<div class="space-y-6">
  <!-- Filter and Controls -->
  <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
    <div class="flex flex-wrap gap-4 items-center">
      <div class="flex-1 min-w-64">
        <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">Filter by name</Label>
        <Input
          bind:value={flpc.filterCond}
          placeholder="Search files and folders..."
          class="w-full"
        />
      </div>
      
      <div>
        <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">Columns</Label>
        <select
          bind:value={flpc.columnCount}
          class="w-32 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="1">1 Column</option>
          <option value="2">2 Columns</option>
          <option value="3">3 Columns</option>
          <option value="4">4 Columns</option>
          <option value="5">5 Columns</option>
        </select>
      </div>
    </div>
  </div>

  <!-- Folders Section -->
  <Toggleable visible={flpc.folderListFiltered().length > 0}>
    {#snippet children()}
      <div class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-xl border border-blue-200 dark:border-blue-800">
        <div class="p-4 border-b border-blue-200 dark:border-blue-700">
          <div class="flex items-center justify-between">
            <h2 class="text-xl font-bold text-blue-900 dark:text-blue-100 flex items-center">
              <div class="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center mr-3">
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                </svg>
              </div>
              Folders
              <span class="ml-2 px-2 py-1 bg-blue-500 text-white text-sm rounded-full">{flpc.folderListFiltered().length}</span>
            </h2>
            <Button
              variant="ghost"
              onclick={() => flpc.showFolders = !flpc.showFolders}
              class="hover:bg-blue-100 dark:hover:bg-blue-800 transition-colors"
              size="sm"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={flpc.showFolders ? "M19 9l-7 7-7-7" : "M9 5l7 7-7 7"}></path>
              </svg>
            </Button>
          </div>
        </div>
        
        <Toggleable visible={flpc.showFolders}>
          {#snippet children()}
            <div class="p-6">
              <div class="grid gap-4" style="grid-template-columns: repeat({flpc.columnCount}, minmax(0, 1fr));">
                <FileCardList
                  list={flpc.folderListFiltered()}
                  callback={callbacks.folderCallback}
                  isFile={false}
                  showSize={false}
                  columnCount={flpc.columnCount}
                />
              </div>
            </div>
          {/snippet}
        </Toggleable>
      </div>
    {/snippet}
  </Toggleable>

  <!-- Files Section -->
  <Toggleable visible={flpc.fileListFiltered().length > 0}>
    {#snippet children()}
      <div class="bg-gradient-to-r from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 rounded-xl border border-green-200 dark:border-green-800">
        <div class="p-4 border-b border-green-200 dark:border-green-700">
          <div class="flex items-center justify-between">
            <h2 class="text-xl font-bold text-green-900 dark:text-green-100 flex items-center">
              <div class="w-8 h-8 bg-green-500 rounded-lg flex items-center justify-center mr-3">
                <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
              </div>
              Files
              <span class="ml-2 px-2 py-1 bg-green-500 text-white text-sm rounded-full">{flpc.fileListFiltered().length}</span>
            </h2>
            <Button
              variant="ghost"
              onclick={() => flpc.showFiles = !flpc.showFiles}
              class="hover:bg-green-100 dark:hover:bg-green-800 transition-colors"
              size="sm"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={flpc.showFiles ? "M19 9l-7 7-7-7" : "M9 5l7 7-7 7"}></path>
              </svg>
            </Button>
          </div>
        </div>
        
        <Toggleable visible={flpc.showFiles}>
          {#snippet children()}
            <div class="p-6">
              <div class="grid gap-4" style="grid-template-columns: repeat({flpc.columnCount}, minmax(0, 1fr));">
                <FileCardList
                  list={flpc.fileListFiltered()}
                  callback={callbacks.fileCallback}
                  download={callbacks.fileDownloadCallback}
                  copy={callbacks.fileCopyCallback}
                  isFile={true}
                  showSize={true}
                  columnCount={flpc.columnCount}
                />
              </div>
            </div>
          {/snippet}
        </Toggleable>
      </div>
    {/snippet}
  </Toggleable>
</div>

<style>
  /* Modern file list panel styles */
</style>
