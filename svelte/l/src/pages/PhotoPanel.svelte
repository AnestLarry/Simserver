<script lang="ts">
  import { Img, Drawer, Listgroup } from "flowbite-svelte";
  import { sineIn } from "svelte/easing";
  import { client } from "../utils";
  export let panel: Panel;

  export let ppc: PhotoPanelConfig = {
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
    <Listgroup
      active
      items={ppc.getListGroup()}
      let:item
      style="word-break:break-all;"
      on:click={(e) => {
        ppc.setShowIndex(e.detail.imgIndex);
      }}
    >
      {item.name}
    </Listgroup>
  {/if}
</Drawer>
<div>
  {#if ppc.imgList != null && ppc.imgList.length > 0}
    {#if panel.photo.photoMode.indexOf("img_2") > -1}
      {#each ppc.imgList as LSItem}
        <Img
          style={panel.photo.sizeRange.isApply
            ? `width:${panel.photo.sizeRange.widthRange}vw;height:${panel.photo.sizeRange.heightRange}vh;`
            : ""}
          src={`${panel.baseUrl}/api/dl/n/${panel.workUrl}` +
            (panel.workUrl.endsWith("/") ? "" : "/") +
            LSItem.Name}
        />
      {/each}
    {:else}
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <div
        on:click={() => {
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
  {/if}
</div>

<style>
</style>
