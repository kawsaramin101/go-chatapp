<script lang="ts">
    import { onDestroy } from "svelte";
    import { goto } from "$app/navigation";

    import { wsStore } from "$lib/stores/ws";

    import Toastify from "toastify-js";
    import "toastify-js/src/toastify.css";

    let username = "";
    let users: string[] = [];
    let chatName = "";

    function createChat(event: SubmitEvent) {
        event.preventDefault();

        const formData = new FormData(event.target as HTMLFormElement);

        const sendingData = {
            action: "CREATECHAT",
            data: {
                usernames: users,
                chatName: chatName,
                isPrivateChat: false,
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
        console.log(event.data);
        const data = JSON.parse(event.data);

        switch (data["action"]) {
            case "CHECK_IF_USER_EXIST":
                const exists = data["data"]["exists"] as boolean;
                const usernameFound = data["data"]["username"] as string;
                if (exists) {
                    users = [usernameFound, ...users];
                    username = "";
                } else {
                    Toastify({
                        text: `User with username ${username} doesn't exist`,
                        duration: 3000,
                        close: true,
                        gravity: "top", // `top` or `bottom`
                        position: "center", // `left`, `center` or `right`
                        stopOnFocus: true, // Prevents dismissing of toast on hover
                        onClick: function () {}, // Callback after click
                    }).showToast();
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

    function checkUser(event: SubmitEvent) {
        event.preventDefault();
        if (users.includes(username)) {
            Toastify({
                text: "User already added",
                duration: 3000,
                close: true,
                gravity: "top",
                position: "center",
                stopOnFocus: true,
                onClick: function () {},
            }).showToast();
            return;
        }
        const sendingData = {
            action: "CHECK_IF_USER_EXIST",
            data: {
                username: username,
            },
        };
        $wsStore?.send(JSON.stringify(sendingData));
    }
</script>

<div class="container">
    <div class="columns">
        <div class="column is-half">
            <h3 class="title is-3 has-text-centered">Create Group Chat</h3>
            <form on:submit={checkUser}>
                <div class="field">
                    <label class="label" for="username">Username</label>
                    <div class="control is-flex is-align-items-center">
                        <input
                            class="input is-flex-grow-1"
                            type="text"
                            name="username"
                            placeholder="Username"
                            bind:value={username}
                        />
                        <button class="button is-info ml-2" type="submit"
                            >Add</button
                        >
                    </div>
                </div>
            </form>
            <form on:submit={createChat}>
                <div class="field">
                    <label class="label" for="chatName">Chat Name</label>
                    <div class="control is-flex is-align-items-center">
                        <input
                            class="input is-flex-grow-1"
                            id="chatName"
                            name="chatName"
                            placeholder="Chat Name"
                            bind:value={chatName}
                        />
                        <div class="control">
                            <button class="button is-primary ml-2" type="submit"
                                >GO</button
                            >
                        </div>
                    </div>
                </div>
            </form>
        </div>

        <div class="column is-half">
            <h3 class="title is-3 has-text-centered">Added Users</h3>
            <div class="has-background-dark p-3">
                <ol class="list">
                    {#each users as user, index}
                        <li class="list-item">
                            <strong>{user}</strong>
                        </li>
                    {/each}
                </ol>
            </div>
        </div>
    </div>
</div>
