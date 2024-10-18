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
        if (chat.name !== null && chat.name !== "") return chat.name;

        if (chat.users.length == 1 && chat.users[0].username == currentUser)
            return "You";

        return chat.users
            .filter((user) => user.username !== currentUser)
            .map((user) => user.username)
            .join(", ");
    }
</script>

<div class="columns is-centered">
    <div class="column is-half">
        <div class="container">
            <div class="tabs is-toggle is-fullwidth">
                <ul>
                    <li class="is-active"><a href="/">Inbox</a></li>
                    <li><a href="message_requests">Requests</a></li>
                </ul>
            </div>

            <ul class="menu-list">
                {#each $chats as chat}
                    <li class="mt-2">
                        <a
                            href="/chat/{chat.ID}"
                            class="is-block has-background-dark p-3 has-text-light"
                        >
                            <span class="is-size-6">{getChatName(chat)}</span>
                            <p class="has-text-grey-lighter is-size-7 mt-1">
                                Lorem Ipsum is simply dummy text of the prin
                            </p>
                        </a>
                    </li>
                {/each}
            </ul>
        </div>
    </div>
</div>

<style>
    .menu-list a {
        transition: background-color 0.3s ease;
        border-radius: 5px;
    }
    .menu-list a:hover {
        background-color: #363636 !important;
    }
</style>
