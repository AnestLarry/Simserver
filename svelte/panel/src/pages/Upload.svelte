<script lang="ts">
  let selectedFile: FileList = $state();
  let progress = $state(0);

  function handleSubmit(e: Event) {
    e.preventDefault();
    if (selectedFile === null || selectedFile.length === 0) {
      alert("Please select a file.");
      return;
    }
    const formData = new FormData();
    formData.append("file", selectedFile[0]);
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/upload/");
    console.log(selectedFile[0].name);
    xhr.setRequestHeader("x-file-name", encodeBase64(selectedFile[0].name));
    xhr.upload.addEventListener("progress", (e) => {
      progress = (e.loaded / e.total) * 100;
    });
    xhr.onloadend = (e) => {
      console.log(e);
      alert("Upload completed.")
    }
    xhr.onabort = (e) => {
      console.error(e);
      alert("Upload aborted.");
    }
    xhr.onerror = (e) => {
      console.error(e);
      alert("Upload failed.");
    }
    xhr.send(formData);
  }
  function encodeBase64(str: string): string {
    return btoa(
      encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, (match, p1) =>
        // @ts-ignore
        String.fromCharCode("0x" + p1)
      )
    );
  }
</script>

<div class="max-w-4xl mx-auto space-y-8">
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6 flex items-center">
      <svg class="w-6 h-6 mr-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
      </svg>
      File Upload
    </h2>
    
    <form onsubmit={handleSubmit} class="space-y-6">
      <div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-8 text-center hover:border-blue-400 transition-colors">
        <svg class="w-12 h-12 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
        </svg>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2" for="file-upload">
          Choose File to Upload
        </label>
        <input 
          id="file-upload"
          type="file" 
          bind:files={selectedFile} 
          class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 file:cursor-pointer cursor-pointer"
        />
        {#if selectedFile && selectedFile.length > 0}
          <p class="mt-2 text-sm text-green-600 dark:text-green-400">
            Selected: {selectedFile[0].name}
          </p>
        {/if}
      </div>
      
      {#if progress > 0}
        <div class="space-y-2">
          <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400">
            <span>Upload Progress</span>
            <span>{Math.round(progress)}%</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
            <div class="bg-blue-600 h-2 rounded-full transition-all duration-300" style="width: {progress}%"></div>
          </div>
        </div>
      {/if}
      
      <button 
        type="submit" 
        class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-6 rounded-lg transition-colors duration-200 flex items-center justify-center"
        disabled={!selectedFile || selectedFile.length === 0}
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
        </svg>
        Upload File
      </button>
    </form>
  </div>

  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6 flex items-center">
      <svg class="w-6 h-6 mr-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
      </svg>
      Text Upload
    </h2>
    
    <form action="/api/upload/text" method="post" enctype="multipart/form-data" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2" for="text-upload">
          Enter your text content:
        </label>
        <textarea 
          id="text-upload"
          name="text" 
          rows="8"
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white resize-vertical"
          placeholder="Type or paste your text here..."
        ></textarea>
      </div>
      <button 
        type="submit" 
        class="w-full bg-green-600 hover:bg-green-700 text-white font-semibold py-3 px-6 rounded-lg transition-colors duration-200 flex items-center justify-center"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
        </svg>
        Upload Text
      </button>
    </form>
  </div>
</div>

<style>
  /* Modern upload component styles */
</style>
