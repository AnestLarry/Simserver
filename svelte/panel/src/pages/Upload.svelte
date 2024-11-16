<script lang="ts">
  let selectedFile: FileList;
  let progress = 0;

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

<div id="panel">
  <form on:submit={handleSubmit}>
    Files: <input type="file" bind:files={selectedFile} /><br /><br />
    <span>Progress: </span><progress value={progress} max="100"></progress><br
    />
    <input type="submit" value="Submit" />
  </form>
  <hr />
  <p>This is a textarea you will upload.</p>
  <form action="/api/upload/text" method="post" enctype="multipart/form-data">
    <textarea name="text"></textarea><br /><br />
    <input type="submit" value="Submit" />
  </form>
</div>

<style>
  #panel {
    width: 45%;
    font-size: 1.25rem;
    margin: 15vh auto;
    background-color: #f2f2f2;
    border-radius: 10px;
    padding: 20px;
    box-shadow: 0px 0px 10px #888888;
  }
  #panel input[type="file"],
  #panel input[type="submit"] {
    padding: 5px;
    display: inline-block;
    border: 1px solid #f44336;
    width: 300px;
    height: 50px;
    border-radius: 5px;
    font-size: 1rem;
    font-weight: bold;
    color: #666666;
    background-color: #f2f2f2;
  }
  #panel input[type="submit"] {
    width: 100%;
    height: 70px;
    color: white;
    background-color: #4caf50;
    border: none;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }
  #panel input[type="submit"]:hover {
    background-color: #3e8e41;
  }
  #panel progress {
    width: 100%;
    height: 20px;
    margin-top: 10px;
    border-radius: 5px;
  }
  #panel textarea {
    width: 100%;
    min-height: 300px;
  }
</style>
