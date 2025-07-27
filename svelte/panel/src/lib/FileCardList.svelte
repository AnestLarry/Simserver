<script lang="ts">
  import Card from "./components/ui/card.svelte";
  import CardContent from "./components/ui/card-content.svelte";
  import Button from "./components/ui/button.svelte";

  let { 
    showSize, 
    list, 
    callback, 
    isFile,
    ondownload,
    oncopy,
    columnCount = "4"
  }: { 
    showSize: Boolean; 
    list: Array<any>; 
    callback: (item: any) => void; 
    isFile: boolean;
    ondownload?: (event: CustomEvent) => void;
    oncopy?: (event: CustomEvent) => void;
    columnCount?: string;
  } = $props();

  // Determine if we should show text based on column count
  let showButtonText = $derived(parseInt(columnCount) <= 3);

  function size(n: number): string {
    let r = [n.toFixed(2), "MB"];
    if (n < 1 || n > 1024) {
      r =
        n < 1 ? [(n * 1024).toFixed(2), " KB"] : [(n / 1024).toFixed(2), "GB"];
    }
    return r.join(" ");
  }
  function download(item: any) {
    ondownload?.(new CustomEvent("download", { detail: item }));
  }
  function copy(item: any) {
    oncopy?.(new CustomEvent("copy", { detail: item }));
  }
</script>

{#each list as LSItem}
  <Card class="hover:shadow-xl hover:scale-[1.02] transition-all duration-200 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 rounded-lg overflow-hidden">
    <CardContent 
      class="p-4 {isFile ? '' : 'cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-750'} transition-colors duration-150"
      onclick={() => {
        if (!isFile) callback(LSItem);
      }}
      role={isFile ? undefined : "button"}
      tabindex={isFile ? undefined : 0}
      onkeydown={isFile ? undefined : (e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          callback(LSItem);
        }
      }}
    >
        <div class="flex items-start justify-between mb-3">
          <div class="flex items-center space-x-3">
            {#if isFile}
              <div class="flex-shrink-0 w-8 h-8 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                <svg class="w-4 h-4 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
              </div>
            {:else}
              <div class="flex-shrink-0 w-8 h-8 bg-yellow-100 dark:bg-yellow-900/30 rounded-lg flex items-center justify-center">
                <svg class="w-4 h-4 text-yellow-600 dark:text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                </svg>
              </div>
            {/if}
            <div class="min-w-0 flex-1">
              <h5 class="font-medium text-gray-900 dark:text-white text-sm break-all leading-tight">
                {LSItem.Name}
              </h5>
            </div>
          </div>
        </div>
      
      <div class="text-xs text-gray-500 dark:text-gray-400 mb-3 space-y-1">
        <div class="flex items-center space-x-1.5">
          <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span class="truncate">{new Date(LSItem.ModTime).toLocaleDateString()}</span>
        </div>
        {#if showSize}
          <div class="flex items-center space-x-1.5">
            <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7c0 2.21-3.582 4-8 4s-8-1.79-8-4z"></path>
            </svg>
            <span class="font-medium text-blue-600 dark:text-blue-400 truncate">{size(LSItem.Size)}</span>
          </div>
        {/if}
      </div>
      
      {#if isFile}
        <div class="flex space-x-1 mt-3 pt-3 border-t border-gray-100 dark:border-gray-700">
          <Button 
            size="sm" 
            variant="default"
            onclick={(e) => {
              e.stopPropagation();
              callback(LSItem);
            }}
            class="flex-1 hover:bg-blue-700 transition-colors text-xs py-1.5 min-w-0"
            title="View file"
          >
            <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
            </svg>
            {#if showButtonText}
              <span class="ml-1 truncate">View</span>
            {/if}
          </Button>
          <Button 
            size="sm" 
            variant="outline"
            onclick={(e) => {
              e.stopPropagation();
              download(LSItem);
            }}
            class="flex-1 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors text-xs py-1.5 min-w-0"
            title="Download file"
          >
            <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            {#if showButtonText}
              <span class="ml-1 truncate">Download</span>
            {/if}
          </Button>
          <Button 
            size="sm" 
            variant="ghost"
            onclick={(e) => {
              e.stopPropagation();
              copy(LSItem);
            }}
            class="flex-1 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-xs py-1.5 min-w-0"
            title="Copy link"
          >
            <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
            </svg>
            {#if showButtonText}
              <span class="ml-1 truncate">Copy</span>
            {/if}
          </Button>
        </div>
      {/if}
    </CardContent>
  </Card>
{/each}

<style>
  /* Custom styles for file cards */
</style>
