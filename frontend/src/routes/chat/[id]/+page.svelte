<script lang="ts">
    import { getContext, onMount } from "svelte";
    import { page } from "$app/stores";
    import { websocket } from "$lib/stores/ws";

    type Message = {
        message: string;
        from: string;
    };

    let messages: Message[] = [];
    let message: string = "";
    let chatContainer;
    let connection: WebSocket;

    $: chatId = Number($page.params.id);

    onMount(() => {
        connection = websocket.get();

        connection.onmessage = function (event) {
            console.log("run");
            const data = JSON.parse(event.data);
            switch (data["action"]) {
                case "MESSAGE":
                    const newMessage: Message = {
                        message: data["data"]["message"],
                        from: data["data"]["from"],
                    };
                    messages = [newMessage, ...messages];
                    break;
                default:
                    break;
            }
        };
    });

    function handleSendMessage(event: SubmitEvent) {
        connection = websocket.get();
        event.preventDefault();
        const sendingData = {
            action: "MESSAGE",
            data: {
                chatId: chatId,
                message: message,
            },
        };
        connection.send(JSON.stringify(sendingData));
    }
</script>

<div class="chat-container" bind:this={chatContainer}>
    <!-- Display messages in the natural order (older at the top, newer at the bottom) -->
    {#each messages as msg}
        <div class="message">
            <strong>{msg.from}:</strong>
            {msg.message}
        </div>
    {/each}
</div>

<form id="form" on:submit={handleSendMessage}>
    <input id="input" autocomplete="off" bind:value={message} /><button
        >Send</button
    >
</form>

<style>
    :global(body) {
        margin: 0;
        padding-bottom: 3rem;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
            Helvetica, Arial, sans-serif;
    }

    #form {
        background: rgba(0, 0, 0, 0.15);
        padding: 0.25rem;
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        display: flex;
        height: 3rem;
        box-sizing: border-box;
        backdrop-filter: blur(10px);
    }
    #input {
        border: none;
        padding: 0 1rem;
        flex-grow: 1;
        border-radius: 2rem;
        margin: 0.25rem;
    }
    #input:focus {
        outline: none;
    }
    #form > button {
        background: #333;
        border: none;
        padding: 0 1rem;
        margin: 0.25rem;
        border-radius: 3px;
        outline: none;
        color: #fff;
    }

    .chat-container {
        max-height: 300px; /* Set a fixed height */
        overflow-y: auto; /* Enable vertical scrolling */
        border: 1px solid #ccc;
        padding: 10px;
        background-color: #f9f9f9;
    }

    .message {
        margin-bottom: 10px;
        padding: 8px;
        background-color: #e1e4ea;
        border-radius: 5px;
    }
</style>
