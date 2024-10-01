<script lang="ts">
    import { getContext, onMount } from "svelte";
    import { chats } from "$lib/stores/chats";
    import { wsStore } from "$lib/stores/ws";
    import { goto } from "$app/navigation";

    const currentUser = localStorage.getItem("username");

    $: if ($wsStore) {
        console.log("WebSocket is ready in this component");
        // Use $wsStore here
    }

    onMount(() => {
        // const connection = websocket.get();

        $wsStore!.onmessage = (event: MessageEvent) => {
            const data = JSON.parse(event.data);

            switch (data["action"]) {
                case "CHAT_CREATED":
                    alert("Chat created. Redirecting");
                    setTimeout(() => {}, 3000);
                    goto(`/chat/${data["data"]["chatId"]}`);
                    break;
                default:
                    break;
            }
        };
    });
</script>

<a href="/login">Login</a>
<a href="/create_chat">Create Chat</a>

<ul>
    {#each $chats as chat}
        <li>
            <a href="/chat/{chat.ID}">
                {#if chat.name}
                    {chat.name}
                {:else}
                    {chat.users
                        .filter((user) => user.username !== currentUser)
                        .map((user) => user.username)
                        .join(", ")}
                {/if}
            </a>
        </li>
    {/each}
</ul>
