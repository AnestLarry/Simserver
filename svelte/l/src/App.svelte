<script lang="ts">
  import PhotoPanel from "./lib/PhotoPanel.svelte";
  import FileListPanel from "./lib/FileListPanel.svelte";
  import {
    Input,
    Card,
    Select,
    Label,
    MultiSelect,
    Button,
    Radio,
  } from "flowbite-svelte";
  let ppc: PhotoPanelConfig;
  let panel: Panel = {
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
  };
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

<div id="list">
  <h1>Items:</h1>
  <Input
    bind:value={panel.workUrl}
    on:mouseleave={() => {
      updateWorkUrl(true);
    }}
    on:blur={() => {
      updateWorkUrl(true);
    }}
  />
  <div class="config">
    <Card class="card">
      <Button
        on:click={() => {
          panel.popUrlStack();
        }}
      >
        Up To Prev Folder
      </Button>
      <br />
      <Button
        on:click={() => {
          window.open(panel.baseUrl + "/api/dl/zip/" + panel.workUrl, "_blank");
        }}
      >
        pack this folder
      </Button>
    </Card>
    <Card class="card">
      <Label>
        Page Mode: {panel.pageMode}
        <Radio name="pageMode" bind:group={panel.pageMode} value="List">
          List
        </Radio>
        <Radio name="pageMode" bind:group={panel.pageMode} value="Photo">
          Photo
        </Radio>
      </Label>
    </Card>
    <Card class="card">
      <Label>
        SortBy
        <Select
          items={[
            { value: "NameOrder", name: "NameOrder" },
            { value: "NameReverse", name: "NameReverse" },
            { value: "TimeOrder", name: "TimeOrder" },
            { value: "TimeReverse", name: "TimeReverse" },
          ]}
          bind:value={panel.sortedBy}
          on:change={() => {
            updateWorkUrl(false);
          }}
        />
      </Label>
    </Card>
    {#if panel.pageMode === "Photo"}
      <Card class="card">
        Photo Mode<br />
        <MultiSelect
          style="width: 100%;"
          items={[
            { value: "img_2", name: "multiple" },
            { value: "doubleFolder", name: "double folder" },
          ]}
          bind:value={panel.photo.photoMode}
        />
        <br />
        <Button on:click={()=>{ppc.showDrawer();}}>
          Show Image List
        </Button>
      </Card>
      <form>
        <Card class="card">
          <Label>
            Photo Width
            <input
              type="range"
              min="50"
              max="150"
              bind:value={panel.photo.sizeRange.widthRange}
              on:mouseleave={() => panel.photo.sizeRange.updateRange()}
            />
            {panel.photo.sizeRange.widthRange}
          </Label>
          <Label>
            Photo Height
            <input
              type="range"
              min="50"
              max="150"
              bind:value={panel.photo.sizeRange.heightRange}
              on:mouseleave={() => panel.photo.sizeRange.updateRange()}
            />
            {panel.photo.sizeRange.heightRange}
          </Label><br />
          <Button
            color="green"
            on:click={() => {
              panel.photo.sizeRange.widthRange = 100;
              panel.photo.sizeRange.heightRange = 100;
            }}
          >
            Reset
          </Button>
          &emsp;
          <Button
            color="light"
            on:click={() => {
              panel.photo.sizeRange.isApply = !panel.photo.sizeRange.isApply;
            }}
          >
            Toggle Effective
          </Button>
        </Card>
      </form>
    {/if}
  </div>
  <div style="clear: both"></div>
  <div class="subPanel">
    {#if panel.pageMode === "List"}
      <FileListPanel {panel} />
    {:else if panel.pageMode === "Photo"}
      <PhotoPanel {panel} bind:ppc />
    {/if}
  </div>
</div>

<style>
  #list {
    margin: 2em 1em 2em 1em;
  }
  .subPanel {
    margin: 1em;
  }
  .config {
    margin: 2em 1em 1em 1em;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }
  .card {
    margin: 1em;
    padding: 1em;
    border: 1px solid black;
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
