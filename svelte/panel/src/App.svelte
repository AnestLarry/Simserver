<script lang="ts">
  import PhotoPanel from "./pages/PhotoPanel.svelte";
  import FileListPanel from "./pages/FileListPanel.svelte";
  import ChatBoard from "./pages/ChatBoard.svelte";
  import Toggleable from "./lib/Toggleable.svelte";
  import Button from "$lib/components/ui/button.svelte";
  import Card from "$lib/components/ui/card.svelte";
  import CardContent from "$lib/components/ui/card-content.svelte";
  import CardHeader from "$lib/components/ui/card-header.svelte";
  import Input from "$lib/components/ui/input.svelte";
  import Label from "$lib/components/ui/label.svelte";
  import Select from "$lib/components/ui/select.svelte";
  import Upload from "./pages/Upload.svelte";
  
  let ppc: PhotoPanelConfig = $state();
  let panel: Panel = $state({
    baseUrl: window.location.protocol + "//" + window.location.host,
    workUrl: "/",
    pageMode: "List",
    sortedBy: "NameOrder",
    photo: {
      photoMode: [],
      sizeRange: {
        isApply: false,
        widthRange: 100,
        heightRange: 100,
        updateRange: () => {
          panel.photo.sizeRange = panel.photo.sizeRange;
        },
      },
    },
    workUrlListening: [],
    hiddenPanel: true,
    pushUrlStack: (x: string) => {
      var urlStack_localStorage = localStorage.getItem("urlStack");
      if (urlStack_localStorage != null) {
        var urlStack: string[] = JSON.parse(urlStack_localStorage);
        if (urlStack.length > 0 && urlStack[0] != "/") {
          urlStack.push(x);
        } else {
          urlStack = ["/" + x];
        }
        panel.workUrl = urlStack.join("/");
        updateWorkUrl(true);
      }
    },
    popUrlStack: () => {
      var urlStack_localStorage = localStorage.getItem("urlStack");
      if (urlStack_localStorage != null) {
        var urlStack: string[] = JSON.parse(urlStack_localStorage);
        urlStack.pop();
        if (urlStack.length > 1) {
          panel.workUrl = urlStack.join("/");
        } else if (urlStack.length == 1) {
          panel.workUrl = urlStack[0];
        } else {
          panel.workUrl = "/";
        }
        updateWorkUrl(true);
      }
    },
  });
  function updateWorkUrl(force: boolean) {
    var value = panel.workUrl.replaceAll("\\", "/");
    value = value.replace(/\/{2,}/g, "/");
    var urlStack_localStorage = localStorage.getItem("urlStack");
    if (urlStack_localStorage != null) {
      if (!force && value === JSON.parse(urlStack_localStorage).join("/")) {
        panel.workUrlListening.forEach((x) => {
          x(force);
        });
        return null;
      }
    }
    var url = value.substring(value.indexOf("/") + 1);
    var urlStack = url.split("/");
    if (!urlStack[0].startsWith("/")) {
      urlStack[0] = "/" + urlStack[0];
    }
    localStorage.setItem("urlStack", JSON.stringify(urlStack));
    panel.workUrlListening.forEach((x) => {
      x(force);
    });
    panel.workUrl = panel.workUrl;
    return null;
  }
  function init() {
    var urlStack_localStorage = localStorage.getItem("urlStack");
    var urlStack = null;
    if (urlStack_localStorage === null) {
      urlStack_localStorage = "";
    } else {
      urlStack = JSON.parse(urlStack_localStorage);
    }
    if (urlStack !== null && Object.keys(urlStack).length > 0) {
      panel.workUrl = urlStack.join("/");
    } else {
      localStorage.setItem("urlStack", JSON.stringify(["/"]));
      panel.workUrl = "/";
    }
  }
  init();
</script>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900">
  <div class="container mx-auto px-4 py-6">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-6">File Manager</h1>
      <div class="flex gap-4 mb-6">
        <div class="flex-1">
          <Input
            bind:value={panel.workUrl}
            onmouseleave={() => {
              updateWorkUrl(true);
            }}
            onblur={() => {
              updateWorkUrl(true);
            }}
            placeholder="Enter path..."
            class="w-full"
          />
        </div>
        <Button
          variant="default"
          onclick={() => {
            panel.hiddenPanel = !panel.hiddenPanel;
            console.log(panel.hiddenPanel);
          }}
          class="hover:bg-blue-700 transition-colors px-3 py-2"
          size="sm"
        >
          <svg class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={panel.hiddenPanel ? "M19 9l-7 7-7-7" : "M9 5l7 7-7 7"}></path>
          </svg>
          {panel.hiddenPanel ? 'Hide' : 'Show'} Panel
        </Button>
      </div>
      <Toggleable visible={panel.hiddenPanel}>
        {#snippet children()}
          <div class="flex flex-wrap gap-4 mb-6">
            <div class="flex gap-2">
              <Button
                variant="outline"
                onclick={() => {
                  panel.popUrlStack();
                }}
                class="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                size="sm"
              >
                <svg class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
                </svg>
                Back
              </Button>
              <Button
                variant="default"
                onclick={() => {
                  window.open(
                    panel.baseUrl + "/api/dl/zip/" + panel.workUrl,
                    "_blank"
                  );
                }}
                class="hover:bg-blue-700 transition-colors"
                size="sm"
              >
                <svg class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                Download
              </Button>
            </div>
            
            <div class="flex items-center gap-2">
              <Label class="text-sm font-medium text-gray-900 dark:text-white whitespace-nowrap">
                Page Mode:
              </Label>
              <Select
                items={[
                  { value: "List", name: "📁 List View" },
                  { value: "Photo", name: "🖼️ Photo Gallery" },
                  { value: "ChatBoard", name: "💬 Chat Board" },
                  { value: "Upload", name: "📤 Upload Files" },
                ]}
                bind:value={panel.pageMode}
                class="min-w-40"
              />
            </div>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">
            {#if panel.pageMode === "List"}
              <Card class="p-4">
                <CardContent class="p-0">
                  <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">
                    Sort By
                  </Label>
                  <Select
                    items={[
                      { value: "NameOrder", name: "Name (A-Z)" },
                      { value: "NameReverse", name: "Name (Z-A)" },
                      { value: "TimeOrder", name: "Date (Oldest)" },
                      { value: "TimeReverse", name: "Date (Newest)" },
                    ]}
                    bind:value={panel.sortedBy}
                    onchange={() => {
                      updateWorkUrl(false);
                    }}
                    class="w-full"
                  />
                </CardContent>
              </Card>
            {:else if panel.pageMode === "Photo"}
              <Card class="p-4 md:col-span-2 lg:col-span-1">
                <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">
                  Photo Mode
                </Label>
                <div class="space-y-2 mb-4">
                  <label class="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={panel.photo.photoMode.includes("img_2")}
                      onchange={(e) => {
                        const target = e.target as HTMLInputElement;
                        if (target.checked) {
                          if (!panel.photo.photoMode.includes("img_2")) {
                            panel.photo.photoMode = [...panel.photo.photoMode, "img_2"];
                          }
                        } else {
                          panel.photo.photoMode = panel.photo.photoMode.filter(x => x !== "img_2");
                        }
                      }}
                      class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span class="text-sm text-gray-700 dark:text-gray-300">Show All Images</span>
                  </label>
                  <label class="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={panel.photo.photoMode.includes("doubleFolder")}
                      onchange={(e) => {
                        const target = e.target as HTMLInputElement;
                        if (target.checked) {
                          if (!panel.photo.photoMode.includes("doubleFolder")) {
                            panel.photo.photoMode = [...panel.photo.photoMode, "doubleFolder"];
                          }
                        } else {
                          panel.photo.photoMode = panel.photo.photoMode.filter(x => x !== "doubleFolder");
                        }
                      }}
                      class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                    />
                    <span class="text-sm text-gray-700 dark:text-gray-300">Include Subfolders</span>
                  </label>
                </div>
                <Button
                  variant="default"
                  onclick={() => {
                    ppc?.showDrawer();
                  }}
                  class="w-full hover:bg-blue-700 transition-colors"
                  size="sm"
                >
                  <svg class="w-3 h-3 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
                  </svg>
                  Image List
                </Button>
              </Card>
              
              <Card class="p-4 md:col-span-2">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Image Size Controls</h3>
                <div class="space-y-4">
                  <div>
                    <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">
                      Width: {panel.photo.sizeRange.widthRange}%
                    </Label>
                    <input
                      type="range"
                      min="50"
                      max="150"
                      bind:value={panel.photo.sizeRange.widthRange}
                      onmouseleave={() => panel.photo.sizeRange.updateRange()}
                      class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
                    />
                  </div>
                  <div>
                    <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">
                      Height: {panel.photo.sizeRange.heightRange}%
                    </Label>
                    <input
                      type="range"
                      min="50"
                      max="150"
                      bind:value={panel.photo.sizeRange.heightRange}
                      onmouseleave={() => panel.photo.sizeRange.updateRange()}
                      class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
                    />
                  </div>
                  <div class="flex gap-2">
                    <Button
                      variant="secondary"
                      onclick={() => {
                        panel.photo.sizeRange.widthRange = 100;
                        panel.photo.sizeRange.heightRange = 100;
                      }}
                      class="flex-1 hover:bg-green-700 transition-colors"
                      size="sm"
                    >
                      Reset
                    </Button>
                    <Button
                      variant={panel.photo.sizeRange.isApply ? "default" : "outline"}
                      onclick={() => {
                        panel.photo.sizeRange.isApply = !panel.photo.sizeRange.isApply;
                      }}
                      class="flex-1 transition-colors {panel.photo.sizeRange.isApply ? 'hover:bg-blue-700' : 'hover:bg-gray-100 dark:hover:bg-gray-700'}"
                      size="sm"
                    >
                      {panel.photo.sizeRange.isApply ? "Applied" : "Apply"}
                    </Button>
                  </div>
                </div>
              </Card>
            {/if}
          </div>
        {/snippet}
      </Toggleable>
    </div>
    
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
      {#if panel.pageMode === "List"}
        <FileListPanel {panel} />
      {:else if panel.pageMode === "Photo"}
        <PhotoPanel {panel} bind:ppc />
      {:else if panel.pageMode === "ChatBoard"}
        <ChatBoard />
      {:else if panel.pageMode === "Upload"}
        <Upload />
      {/if}
    </div>
  </div>
</div>

<style>
  /* Custom range slider styling */
  input[type="range"]::-webkit-slider-thumb {
    appearance: none;
    height: 20px;
    width: 20px;
    border-radius: 50%;
    background: #3b82f6;
    cursor: pointer;
    border: 2px solid #ffffff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  input[type="range"]::-moz-range-thumb {
    height: 20px;
    width: 20px;
    border-radius: 50%;
    background: #3b82f6;
    cursor: pointer;
    border: 2px solid #ffffff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  input[type="range"]::-webkit-slider-track {
    height: 8px;
    border-radius: 4px;
    background: #e5e7eb;
  }

  input[type="range"]::-moz-range-track {
    height: 8px;
    border-radius: 4px;
    background: #e5e7eb;
  }
</style>