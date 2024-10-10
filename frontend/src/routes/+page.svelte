<script lang="ts">
    import { getContext, onMount } from "svelte";
    import { chats } from "$lib/stores/chats";
    import { wsStore } from "$lib/stores/ws";
    import type { Chat } from "$lib/models";

    const currentUser = localStorage.getItem("username");

    $: if ($wsStore) {
        console.log("WebSocket is ready in this component");
        // Use $wsStore here
    }

    function getChatName(chat: Chat) {
        if (chat.name !== "") return chat.name;

        if (chat.users.length == 1 && chat.users[0].username == currentUser)
            return "You";

        return chat.users
            .filter((user) => user.username !== currentUser)
            .map((user) => user.username)
            .join(", ");
    }
</script>

<div class="section">
    <div class="container has-text-centered">
        <!-- Centered Login and Create Chat links -->
        <div class="buttons is-centered">
            <a class="button is-primary" href="/login">Login</a>
        </div>

        <ul>
            {#each $chats as chat}
                <li class="mt-4">
                    <a class="has-text-weight-bold" href="/chat/{chat.ID}">
                        {getChatName(chat)}
                    </a>
                </li>
            {/each}
        </ul>
    </div>
</div>
