<script lang="ts">
  import Button from "../lib/components/ui/button.svelte";
  import Drawer from "../lib/components/ui/drawer.svelte";
  import List from "../lib/components/ui/list.svelte";
  import { sineIn } from "svelte/easing";
  import { client } from "../utils";
  
  let { panel, ppc = $bindable() }: { panel: Panel; ppc: PhotoPanelConfig } = $props();

  // Initialize ppc with $state if not provided
  if (!ppc) {
    ppc = {
    imgList: [],
    showIndex: 0,
    listenIndex: -1,
    drawerHidden: true,
    showDrawer: () => {
      ppc.drawerHidden = false;
      setTimeout(() => {
        var cur = document.querySelector("button[aria-current='true']");
        if (cur !== null) {
          cur.focus();
        }
      }, 50);
    },
    fresh: (f: boolean) => {
      if (!f) {
        return;
      }
      ppc.imgList = [];
      if (panel.photo.photoMode.indexOf("doubleFolder") > -1) {
        var folders = client.FolderList(
          panel.baseUrl + "/api/dl/ls/" + panel.workUrl
        );
        folders.forEach((f) => {
          ppc.imgList = ppc.imgList.concat(
            client.ImgList(
              `${panel.baseUrl}/api/dl/ls/${panel.workUrl}` +
                (panel.workUrl.endsWith("/") ? "" : "/") +
                f.Name,
              f.Name + (f.Name.endsWith("/") ? "" : "/")
            )
          );
        });
      } else {
        ppc.imgList = client.ImgList(
          `${panel.baseUrl}/api/dl/ls/${panel.workUrl}`,
          ""
        );
      }
      ppc.setShowIndex(0);
    },
    setShowIndex: (i: number) => {
      if (i < ppc.imgList.length) {
        ppc.showIndex = i;
      } else {
        ppc.showIndex = i % ppc.imgList.length;
      }
    },
    getListGroup: () => {
      let res: any[] = [];
      let imgs = ppc.imgList.map((x) => x);
      imgs.sort(client.sortFunction("NameLenOrder"));
      imgs.forEach((x) => {
        let i = ppc.imgList.indexOf(x);
        res.push({
          name: x.Name,
          current: i === ppc.showIndex,
          imgIndex: i,
        });
      });
        return res;
      },
    };
  }

  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn,
  };
  
  function init() {
    ppc.fresh(true);
    ppc.listenIndex = panel.workUrlListening.push(ppc.fresh);
  }
  init();
</script>

<Drawer transitionType="fly" {transitionParams} bind:hidden={ppc.drawerHidden}>
  {#if ppc.imgList.length > 0}
    <List
      active
      items={ppc.getListGroup()}
      class="break-all"
      onclick={(e) => {
        ppc.setShowIndex(e.detail.imgIndex);
      }}
    >
      {#snippet children({ item })}
        {item.name}
      {/snippet}
    </List>
  {/if}
</Drawer>
<div>
  {#if ppc.imgList != null && ppc.imgList.length > 0}
    {#if panel.photo.photoMode.indexOf("img_2") > -1}
      {#each ppc.imgList as LSItem}
        <img
          class="rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200"
          style={panel.photo.sizeRange.isApply
            ? `width:${panel.photo.sizeRange.widthRange}vw;height:${panel.photo.sizeRange.heightRange}vh;`
            : ""}
          src={`${panel.baseUrl}/api/dl/n/${panel.workUrl}` +
            (panel.workUrl.endsWith("/") ? "" : "/") +
            LSItem.Name}
          alt={LSItem.Name}
        />
      {/each}
    {:else}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        onclick={() => {
          ppc.setShowIndex(ppc.showIndex + 1);
        }}
      >
        <img
          style={panel.photo.sizeRange.isApply
            ? `width:${panel.photo.sizeRange.widthRange}vw;height:${panel.photo.sizeRange.heightRange}vh;`
            : ""}
          src={`${panel.baseUrl}/api/dl/n/${panel.workUrl}` +
            (panel.workUrl.endsWith("/") ? "" : "/") +
            ppc.imgList[ppc.showIndex].Name}
          alt={ppc.imgList[ppc.showIndex].Name}
        />
      </div>
    {/if}
  {:else}
    <div class="flex flex-col items-center justify-center py-16 px-4">
      <div class="bg-gradient-to-br from-yellow-50 to-orange-50 dark:from-yellow-900/20 dark:to-orange-900/20 rounded-xl p-8 max-w-lg text-center border border-yellow-200 dark:border-yellow-800 shadow-lg">
        <div class="w-16 h-16 bg-yellow-100 dark:bg-yellow-900/40 rounded-full flex items-center justify-center mx-auto mb-6">
          <svg class="w-8 h-8 text-yellow-600 dark:text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
          </svg>
        </div>
        <h3 class="text-xl font-bold text-yellow-800 dark:text-yellow-200 mb-3">No Images Found</h3>
        <p class="text-yellow-700 dark:text-yellow-300 mb-4 leading-relaxed">No images were found in this directory.</p>
        <div class="bg-yellow-100 dark:bg-yellow-900/30 rounded-lg p-4">
          <p class="text-sm font-medium text-yellow-800 dark:text-yellow-200 mb-2">Possible reasons:</p>
          <ul class="text-sm text-yellow-700 dark:text-yellow-300 space-y-1 text-left">
            <li class="flex items-center">
              <svg class="w-3 h-3 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="2"/>
              </svg>
              The ls feature is disabled in the backend
            </li>
            <li class="flex items-center">
              <svg class="w-3 h-3 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="2"/>
              </svg>
              No image files exist in this location
            </li>
            <li class="flex items-center">
              <svg class="w-3 h-3 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <circle cx="10" cy="10" r="2"/>
              </svg>
              Access permissions may be restricted
            </li>
          </ul>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Photo panel custom styles */
</style>
