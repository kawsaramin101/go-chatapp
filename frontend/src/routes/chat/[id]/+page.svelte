<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import { page } from "$app/stores";
    import { wsStore } from "$lib/stores/ws";

    import { activeChatId } from "$lib/stores/activeChatIdstore";
    import { messageStore } from "$lib/stores/messagesStore";
    import { getAllMessagesFromAChat } from "$lib/storage/messages";
    import type { Message } from "$lib/models";

    let messages: Message[] = [];
    let message: string = "";
    let chatContainer;

    $: chatId = Number($page.params.id);

    onMount(async () => {
        activeChatId.set(chatId);
        messageStore.set(await getAllMessagesFromAChat(chatId));
        console.log(await getAllMessagesFromAChat(chatId));
    });

    function handleIncommingMessage(event: MessageEvent) {
        const data = JSON.parse(event.data);

        // switch (data["action"]) {
        //     case "MESSAGE":
        //         console.log(data["data"]["chatId"] as Number);
        //         if ((data["data"]["chatId"] as Number) === chatId) {
        //             const newMessage: Message = {
        //                 message: data["data"]["message"],
        //                 from: data["data"]["from"],
        //             };
        //             messages = [newMessage, ...messages];
        //         }
        //         break;
        //     default:
        //         break;
        // }
    }

    $: if ($wsStore) {
        $wsStore.addEventListener("message", handleIncommingMessage);
    }

    onDestroy(() => {
        activeChatId.set(null);
        $wsStore?.removeEventListener("message", handleIncommingMessage);
        messageStore.set([]);
    });

    function handleSendMessage(event: SubmitEvent) {
        event.preventDefault();
        const sendingData = {
            action: "MESSAGE",
            data: {
                chatId: chatId,
                message: message,
            },
        };
        $wsStore?.send(JSON.stringify(sendingData));
        message = "";
    }
</script>

<div class="chat-container" bind:this={chatContainer}>
    <!-- Display messages in the natural order (older at the top, newer at the bottom) -->
    {#each $messageStore as msg}
        <div class="message">
            <strong>{msg.from}:</strong>
            {msg.content}
        </div>
    {/each}
</div>

<form id="form" on:submit={handleSendMessage}>
    <input id="input" autocomplete="off" bind:value={message} />
    <button>Send</button>
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
        cursor: pointer;
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
