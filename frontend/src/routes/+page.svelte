<script lang="ts">
    import { getContext, onMount } from "svelte";
    import { chats } from "$lib/stores/chats";
    import { websocket } from "$lib/stores/ws";
    import { goto } from "$app/navigation";

    onMount(() => {
        const connection = websocket.get();

        connection.onmessage = (event: MessageEvent) => {
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

    function addUser(event: SubmitEvent) {
        event.preventDefault();
        const connection = websocket.get();

        console.log(connection);
        const formData = new FormData(event.target as HTMLFormElement);

        const sendingData = {
            action: "CREATECHAT",
            data: {
                username: formData.get("username"),
            },
        };
        connection.send(JSON.stringify(sendingData));
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
