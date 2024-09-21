<script context="module">
</script>

<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { setContext } from "svelte";

    const baseUrl: string = "localhost:8000";

    let currentRoute: string;

    let connection: WebSocket;

    function connectWebSocket() {
        connection = new WebSocket("ws://" + baseUrl + "/ws");

        connection.onopen = function () {
            console.log("WebSocket connection established successfully.");
            const authToken = localStorage.getItem("authToken") || "";
            connection.send(authToken);

            const data = {
                action: "BROADCAST",
                data: {
                    message: "Hello world",
                },
            };

            connection.send(JSON.stringify(data));
        };

        connection.onmessage = function (event) {
            console.log(event.data);
            const data = JSON.parse(event.data);
            if (
                data["action"] === "ERROR_USER_NOT_FOUND" ||
                data["action"] === "ERROR_SERVER_ERROR" ||
                data["action"] === "ERROR_INVALID_PAYLOAD"
            ) {
                alert(data["message"]);
            } else if (data["action"] === "CHAT_CREATED") {
                alert("Chat created");
                setTimeout(() => {}, 3000);
            } else if (data["action"] === "INITIAL_DATA") {
            }
        };

        connection.onclose = function (event) {
            console.log("Websocket connection closed", event);
            let retry: boolean = true;
            if (retry && !event.wasClean) {
                setTimeout(function () {
                    connectWebSocket();
                }, 4000); // Retry after 5 seconds
            }
            // var item = document.createElement("div");
            // item.innerHTML = "<b>Connection closed.</b>";
            // appendLog(item);
        };
    }

    onMount(() => {
        // conn = new WebSocket("ws://" + baseUrl + "/ws");
        connectWebSocket();

        if (currentRoute !== "/login" && currentRoute !== "/signup") {
        }
        console.log("This code runs on all routes");

        document.body.classList.add("js-enabled");
    });

    export function addUser(event: SubmitEvent) {
        event.preventDefault();

        const formData = new FormData(event.target as HTMLFormElement);

        const sendingData = {
            action: "CREATECHAT",
            data: {
                username: formData.get("username"),
            },
        };

        connection.send(JSON.stringify(sendingData));
    }

    setContext("addUser", addUser);

    // Subscribe to the page store to know when the route changes
    $: {
        currentRoute = $page.url.pathname;
        console.log("Current route:", currentRoute);

        // You can add more code here that should run on every route change
    }
</script>

<slot />
