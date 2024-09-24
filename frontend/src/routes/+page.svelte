<script lang="ts">
    import { getContext, onMount } from "svelte";
    import { chats } from "$lib/stores/chats";
    import { websocket } from "$lib/stores/ws";

    let connection: WebSocket;

    onMount(() => {
        connection = websocket.get();
    });

    function addUser(event: SubmitEvent) {
        connection = websocket.get();
        event.preventDefault();

        const formData = new FormData(event.target as HTMLFormElement);

        const sendingData = {
            action: "CREATECHAT",
            data: {
                username: formData.get("username"),
            },
        };
        if (connection !== null) {
            connection.send(JSON.stringify(sendingData));
        }
    }
</script>

<form on:submit={addUser}>
    <input type="text" name="username" placeholder="Username" />
    <button type="submit">GO</button>
</form>
<a href="/login">Login</a>

<ul>
    {#each $chats as chat}
        <li>
            <a href="/chat/{chat.ID}">
                Chat ID: {chat.ID}, Users: {chat.users
                    .map((user) => user.username)
                    .join(", ")}
            </a>
        </li>
    {/each}
</ul>
