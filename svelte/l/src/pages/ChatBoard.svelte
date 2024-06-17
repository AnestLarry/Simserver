<script lang="ts">
  import { Button, Input, Textarea, Card } from "flowbite-svelte";
  import { onMount } from "svelte";

  let ws: WebSocket;
  let username: string = "user" + Math.floor(Math.random() * 10000);
  let message = "";
  let messages: BroadcastMessage[] = [];
  let retryCount = 0;
  let retryTimeout = 1000;
  const maxRetryCount = 3;

  function connect() {
    ws = new WebSocket(
      (window.location.protocol.indexOf("https") > -1 ? "wss://" : "ws://") +
        window.location.host +
        "/api/chatBoard/chat"
    );

    ws.onopen = () => {
      retryCount = 0;
      retryTimeout = 1000;
      console.log("WebSocket connection opened.");
    };

    ws.onmessage = (event) => {
      let bm: BroadcastMessage = JSON.parse(event.data);
      messages = [bm, ...messages];
    };

    ws.onclose = (event) => {
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
  <p>
    Status: {ws != null && ws.readyState === ws.CONNECTING
      ? "CONNECTING"
      : "Fail"}
  </p>
  <Input type="text" bind:value={username} placeholder="username" />
  <Textarea
    type="text"
    bind:value={message}
    on:keydown={handleKeydown}
    placeholder="Type your message here"
  />
  <Button on:click={sendMessage}>Send</Button>  <Button
    on:click={clearMessage}>Clear</Button
  >
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
