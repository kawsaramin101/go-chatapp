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

    const currentUser = localStorage.getItem("username");

    $: chatId = Number($page.params.id);

    onMount(async () => {
        activeChatId.set(chatId);
        messageStore.set(await getAllMessagesFromAChat(chatId));
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

<div class="container">
    <div class="columns is-centered">
        <div class="column is-half">
            <ul id="messages">
                {#each $messageStore as msg}
                    {#if currentUser === msg.from}
                        <li class="message-right">
                            <small>
                                {msg.createdAt.toLocaleString()} - {msg.from}
                            </small>
                            <br />
                            {msg.content}
                        </li>
                    {:else}
                        <li class="message-left">
                            <small
                                >{msg.from} - {msg.createdAt.toLocaleString()}</small
                            >
                            <br />
                            {msg.content}
                        </li>
                    {/if}
                {/each}
            </ul>
            <form id="form" on:submit={handleSendMessage}>
                <input id="input" bind:value={message} autocomplete="off" />
                <button>Send</button>
            </form>
        </div>
    </div>
</div>

<style>
    :global(body) {
        margin: 0;
        padding-bottom: 3rem;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
            Helvetica, Arial, sans-serif;
        background-color: #121212;
        color: #e0e0e0;
    }

    #form {
        background: rgba(0, 0, 0, 1);
        padding: 0.25rem;
        position: fixed;
        bottom: 0;
        left: 0;
        right: 0;
        display: flex;
        height: 3rem;
        box-sizing: border-box;
        max-width: 50%;
        margin-left: auto;
        margin-right: auto;
    }

    #input {
        border: none;
        padding: 0 1rem;
        flex-grow: 1;
        border-radius: 2rem;
        margin: 0.25rem;
        background-color: #3c3c3c;
        color: #e0e0e0;
    }

    #input:focus {
        outline: none;
        background-color: #3c3c3c;
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

    #messages {
        list-style-type: none;
        margin: 0;
        padding: 0;
        max-width: 100%;
        margin-left: auto;
        margin-right: auto;
        margin-top: 2rem;
        padding-bottom: 2rem;
        display: flex;
        flex-direction: column-reverse; /* Invert the order */
        max-height: 85vh;
        overflow-y: auto; /* Enable scrolling */
    }

    .message-left {
        padding: 0.5rem 1rem;
        margin: 5px;
        border-radius: 15px;
        background-color: #0056b3;
        color: white;
        max-width: 70%;
        word-wrap: break-word;
        text-align: left;
        box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.5);
        align-self: flex-start; /* Align to the left */
    }

    .message-right {
        padding: 0.5rem 1rem;
        margin: 5px;
        border-radius: 15px;
        background-color: #1f1f1f;
        color: #e0e0e0;
        max-width: 70%;
        word-wrap: break-word;
        text-align: right;
        box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.5);
        align-self: flex-end; /* Align to the right */
    }
</style>
