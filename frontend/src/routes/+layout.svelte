<script lang="ts">
    import { page } from "$app/stores";
    import { beforeNavigate, goto } from "$app/navigation";

    // import { onMount } from "svelte";
    // import { setContext } from "svelte";
    // import { chats } from "$lib/stores/chats";

    // const baseUrl: string = "localhost:8000";

    // let currentRoute: string;

    // let connection: WebSocket;

    // function connectWebSocket() {
    //     connection = new WebSocket("ws://" + baseUrl + "/ws");

    //     connection.onopen = function () {
    //         console.log("WebSocket connection established successfully.");
    //         const authToken = localStorage.getItem("authToken") || "";
    //         connection.send(authToken);

    //         const data = {
    //             action: "BROADCAST",
    //             data: {
    //                 message: "Hello world",
    //             },
    //         };

    //         connection.send(JSON.stringify(data));
    //     };

    //     connection.onmessage = function (event) {
    //         console.log(event.data);
    //         const data = JSON.parse(event.data);
    //         switch (data["action"]) {
    //             case "ERROR_USER_NOT_FOUND":
    //             case "ERROR_SERVER_ERROR":
    //             case "ERROR_INVALID_PAYLOAD":
    //                 alert(data["message"]);
    //                 break;

    //             case "CHAT_CREATED":
    //                 alert("Chat created");
    //                 setTimeout(() => {}, 3000);
    //                 break;

    //             case "INITIAL_DATA":
    //                 chats.setChats(data["data"]["chats"]);
    //                 break;

    //             case "MESSAGE":

    //             default:
    //                 // Handle any other actions if needed
    //                 break;
    //         }
    //     };

    //     connection.onclose = function (event) {
    //         console.log("Websocket connection closed", event);
    //         let retry: boolean = true;
    //         if (retry && !event.wasClean) {
    //             setTimeout(function () {
    //                 connectWebSocket();
    //             }, 4000); // Retry after 5 seconds
    //         }
    //         // var item = document.createElement("div");
    //         // item.innerHTML = "<b>Connection closed.</b>";
    //         // appendLog(item);
    //     };
    // }

    // onMount(() => {
    //     // conn = new WebSocket("ws://" + baseUrl + "/ws");
    //     connectWebSocket();
    //     setContext("connection", connection);

    //     if (currentRoute !== "/login" && currentRoute !== "/signup") {
    //     }
    //     console.log("This code runs on all routes");
    // });

    // export function addUser(event: SubmitEvent) {
    //     event.preventDefault();

    //     const formData = new FormData(event.target as HTMLFormElement);

    //     const sendingData = {
    //         action: "CREATECHAT",
    //         data: {
    //             username: formData.get("username"),
    //         },
    //     };

    //     connection.send(JSON.stringify(sendingData));
    // }

    // function sendMessage(chatId: number, message: string) {
    //     const sendingData = {
    //         action: "MESSAGE",
    //         data: {
    //             chatId: chatId,
    //             message: message,
    //         },
    //     };
    //     connection.send(JSON.stringify(sendingData));
    // }

    // setContext("addUser", addUser);
    // setContext("sendMessage", sendMessage);

    // // Subscribe to the page store to know when the route changes
    // $: {
    //     currentRoute = $page.url.pathname;
    //     console.log("Current route:", currentRoute);
    // }
    //

    import { onDestroy, onMount } from "svelte";
    import { initializeWebSocket, closeWebSocket } from "$lib/stores/ws";

    onMount(() => {
        initializeWebSocket();
    });

    onDestroy(() => {
        closeWebSocket();
    });

    let currentRoute: string;
    const username = localStorage.getItem("username");
    const authToken = localStorage.getItem("authToken");

    beforeNavigate(({ to, cancel }) => {
        if (to) {
            const route = to.url.pathname;

            const username = localStorage.getItem("username");
            const authToken = localStorage.getItem("authToken");

            if (route !== "/login" && route !== "/signup") {
                if (authToken === null || authToken === "") goto("/login");
            }
        }
    });

    $: {
        currentRoute = $page.url.pathname;
        console.log("Current route:", currentRoute);
        if (currentRoute !== "/login" && currentRoute !== "/signup") {
            if (authToken === null || authToken === "") goto("/login");
        }
    }

    function logout() {
        localStorage.removeItem("authToken");
        localStorage.removeItem("username");
        goto("/login");
    }
</script>

{#if currentRoute !== "/login" && currentRoute !== "/signup"}
    <header class="navbar is-dark is-fixed-top">
        <div class="navbar-brand">
            <h3 class="navbar-item subtitle">{username}</h3>
        </div>
        <div class="navbar-end">
            <div class="navbar-item">
                <a class="button is-link" href="/create_chat">Create Chat</a>
                <button class="button is-light" on:click={logout}>Logout</button
                >
            </div>
        </div>
    </header>
{/if}

<section class="section">
    <div class="container">
        <slot />
    </div>
</section>
