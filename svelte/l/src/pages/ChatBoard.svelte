<script lang="ts">
  import { Button, Input, Textarea, Card } from "flowbite-svelte";
  import axios from "axios";
  import { onMount } from "svelte";

  let ws: WebSocket;
  let username: string = "user" + Math.floor(Math.random() * 10000);
  let message = "";
  let messages: BroadcastMessage[] = [];
  let retryCount = 0;
  let retryTimeout = 1000;
  const maxRetryCount = 3;
  let wsStatus: string = "Disabled";
  function connect() {
    axios
      .get("/api/view/chatBoard/health")
      .then((res) => {
        ws = new WebSocket(
          (window.location.protocol.indexOf("https") > -1
            ? "wss://"
            : "ws://") +
            window.location.host +
            "/api/view/chatBoard/chat"
        );

        ws.onopen = () => {
          retryCount = 0;
          retryTimeout = 1000;
          console.log("WebSocket connection opened.");
          wsStatus = "Connected";
          // update status
          getStatus();
        };

        ws.onmessage = (event) => {
          let bm: BroadcastMessage = JSON.parse(event.data);
          messages = [bm, ...messages];
        };

        ws.onclose = (event) => {
          wsStatus = "Closed";
          console.log("WebSocket connection closed.", event);
          if (retryCount < maxRetryCount) {
            retryCount++;
            retryTimeout *= 2; // 使用指数退避策略增加重连超时时间
            setTimeout(connect, retryTimeout);
          }
        };

        ws.onerror = (error) => {
          console.error("WebSocket error:", error);
          ws.close();
        };
      })
      .catch((err) => {
        wsStatus = "Disabled";
      });
  }

  function getStatus() {
    if (!wsStatus) {
      return "Disable";
    } else if (ws != null && ws.readyState === ws.OPEN) {
      return "Connected";
    } else {
      return "Fail";
    }
  }

  onMount(() => {
    connect();
  });

  function sendMessage() {
    if (message) {
      let um: UserMessage = {
        username: username,
        content: message,
      };
      ws.send(JSON.stringify(um));
      message = "";
    }
  }

  function clearMessage() {
    messages = [];
  }

  function handleKeydown(event) {
    if (event.key === "Enter") {
      sendMessage();
    }
  }
</script>

<div>
  <p style="margin-bottom: 1em;">
    Status: {#if wsStatus === "Connected"}
      <Button color="green" pill>Connected</Button>
    {:else}
      <Button color="red" pill on:click={connect}>{wsStatus}</Button>
    {/if}
  </p>
  <Input type="text" bind:value={username} placeholder="username" disabled={wsStatus !== "Connected"}/>
  <Textarea
    type="text"
    bind:value={message}
    on:keydown={handleKeydown}
    placeholder="Type your message here"
    disabled={wsStatus !== "Connected"}
  />
  <Button on:click={sendMessage} disabled={wsStatus !== "Connected"}>
    Send
  </Button>
  <Button on:click={clearMessage} disabled={wsStatus !== "Connected"}>Clear</Button>
</div>

<ul>
  {#each messages as msg}
    <li>
      {new Date(Number(msg.timestamp) * 1000).toLocaleString()}
      <strong>{msg.username}</strong>({msg.userIp}) -->
      <Textarea bind:value={msg.content} readonly />
    </li>
  {/each}
</ul>

<style>
  ul li {
    padding-bottom: 1em;
    border-bottom: solid black 1px;
  }
</style>
