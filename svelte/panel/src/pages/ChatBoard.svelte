<script lang="ts">
  import Button from "../lib/components/ui/button.svelte";
  import Input from "../lib/components/ui/input.svelte";
  import Card from "../lib/components/ui/card.svelte";
  import Textarea from "../lib/components/ui/textarea.svelte";
  import Label from "../lib/components/ui/label.svelte";
  import axios from "axios";
  import { onMount } from "svelte";

  let ws: WebSocket;
  let username: string = $state("user" + Math.floor(Math.random() * 10000));
  let message = $state("");
  let messages: BroadcastMessage[] = $state([]);
  let retryCount = 0;
  let retryTimeout = 1000;
  const maxRetryCount = 3;
  let wsStatus: string = $state("Disabled");
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

  function handleKeydown(event: any) {
    if (event.key === "Enter") {
      sendMessage();
    }
  }
</script>

<div class="space-y-6">
  <Card class="p-4">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white">Chat Board</h2>
      <div class="flex items-center gap-2">
        <span class="text-sm text-gray-600 dark:text-gray-400">Status:</span>
        {#if wsStatus === "Connected"}
          <Button variant="secondary" class="bg-green-600 hover:bg-green-700 text-white rounded-full px-3 py-1 text-xs">Connected</Button>
        {:else}
          <Button variant="destructive" class="rounded-full px-3 py-1 text-xs" onclick={connect}>{wsStatus}</Button>
        {/if}
      </div>
    </div>
    
    <div class="space-y-4">
      <div>
        <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">Username</Label>
        <Input 
          type="text" 
          bind:value={username} 
          placeholder="Enter your username" 
          disabled={wsStatus !== "Connected"}
          class="w-full"
        />
      </div>
      
      <div>
        <Label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white">Message</Label>
        <Textarea
          bind:value={message}
          onkeydown={handleKeydown}
          placeholder="Type your message here..."
          disabled={wsStatus !== "Connected"}
          class="w-full min-h-24 resize-none"
        />
      </div>
      
      <div class="flex gap-3">
        <Button 
          onclick={sendMessage} 
          disabled={wsStatus !== "Connected"}
          class="flex-1 hover:bg-blue-700 transition-colors"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
          </svg>
          Send
        </Button>
        <Button 
          variant="outline" 
          onclick={clearMessage} 
          disabled={wsStatus !== "Connected"}
          class="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
          </svg>
          Clear
        </Button>
      </div>
    </div>
  </Card>
</div>

{#if messages.length > 0}
  <Card class="p-4 mt-6">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
      </svg>
      Messages ({messages.length})
    </h3>
    <div class="space-y-4 max-h-96 overflow-y-auto">
      {#each messages as msg}
        <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 border border-gray-200 dark:border-gray-600">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
                <span class="text-white text-sm font-medium">{msg.username.charAt(0).toUpperCase()}</span>
              </div>
              <div>
                <span class="font-semibold text-gray-900 dark:text-white">{msg.username}</span>
                <span class="text-xs text-gray-500 dark:text-gray-400 ml-2">({msg.userIp})</span>
              </div>
            </div>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              {new Date(Number(msg.timestamp) * 1000).toLocaleString()}
            </span>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded p-3 border border-gray-200 dark:border-gray-600">
            <p class="text-gray-900 dark:text-white whitespace-pre-wrap break-words">{msg.content}</p>
          </div>
        </div>
      {/each}
    </div>
  </Card>
{/if}
