<script lang="ts">
    import { wsStore } from "$lib/stores/ws";
    import { onMount } from "svelte";

    let username = "";

    function addUser(event: SubmitEvent) {
        event.preventDefault();

        const formData = new FormData(event.target as HTMLFormElement);

        const sendingData = {
            action: "CREATECHAT",
            data: {
                username: formData.get("username"),
            },
        };
        $wsStore?.send(JSON.stringify(sendingData));
    }

    $: if ($wsStore) {
        $wsStore.onmessage = function (event) {
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
                default:
                    break;
            }
        };
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
    <input
        type="text"
        name="username"
        placeholder="Username"
        bind:value={username}
    />
    <button type="button" on:click={checkUser}>Add</button>

    <button type="submit">GO</button>
</form>
