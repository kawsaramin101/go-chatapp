<script lang="ts">
    import { onDestroy } from "svelte";
    import { goto } from "$app/navigation";

    import { wsStore } from "$lib/stores/ws";

    import Toastify from "toastify-js";
    import "toastify-js/src/toastify.css";

    let username = "";

    function addUser(event: SubmitEvent) {
        event.preventDefault();

        const formData = new FormData(event.target as HTMLFormElement);

        const sendingData = {
            action: "CREATECHAT",
            data: {
                usernames: [formData.get("username")],
            },
        };
        $wsStore?.send(JSON.stringify(sendingData));
    }

    $: if ($wsStore) {
        $wsStore.addEventListener("message", handleIncommingMessage);
    }

    onDestroy(() => {
        $wsStore?.removeEventListener("message", handleIncommingMessage);
    });

    function handleIncommingMessage(event: MessageEvent) {
        const data = JSON.parse(event.data);
        switch (data["action"]) {
            case "CHECK_IF_USER_EXIST":
                const exists = data["data"]["exists"] as boolean;
                const username = data["data"]["username"] as string;
                if (exists) {
                    users = [username, ...users];
                } else {
                    alert(`User with username ${username} doesn't exist`);
                }
                break;

            case "CHAT_CREATED":
                Toastify({
                    text: "Chat created. Redirecting",
                    duration: 3000,
                    close: true,
                    gravity: "top", // `top` or `bottom`
                    position: "center", // `left`, `center` or `right`
                    stopOnFocus: true, // Prevents dismissing of toast on hover
                    onClick: function () {}, // Callback after click
                }).showToast();
                setTimeout(() => {}, 3000);

                goto(`/chat/${data["data"]["chatId"]}`);
                break;

            default:
                break;
        }
    }

    function checkUser() {
        const sendingData = {
            action: "CHECK_IF_USER_EXIST",
            data: {
                username: username,
            },
        };
        $wsStore?.send(JSON.stringify(sendingData));
    }

    function anotherUser() {}

    let users: string[] = [];
</script>

<form on:submit={addUser}>
    <input type="text" name="username" placeholder="Username" required />
    <button type="submit">GO</button>
</form>

<h3>Create Group Chat</h3>
<form>
    <input id="chatName" name="chatName" value="" />
    <br />
    <input
        type="text"
        name="username"
        placeholder="Username"
        bind:value={username}
    />

    <button type="button" on:click={checkUser}>Add</button>
    <br />
    <ol>
        {#each users as user, index}
            <li>{user}</li>
        {/each}
    </ol>
    <br />

    <button type="submit">GO</button>
</form>
