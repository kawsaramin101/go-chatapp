<script lang="ts">
    import { getContext, onMount } from "svelte";
    import { chats } from "$lib/stores/chats";
    import { wsStore } from "$lib/stores/ws";

    const currentUser = localStorage.getItem("username");

    $: if ($wsStore) {
        console.log("WebSocket is ready in this component");
        // Use $wsStore here
    }
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
