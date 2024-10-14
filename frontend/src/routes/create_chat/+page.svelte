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
                isPrivateChat: true,
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
</script>

<form on:submit={addUser}>
    <input type="text" name="username" placeholder="Username" required />
    <button type="submit">GO</button>
</form>
