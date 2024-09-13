<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";

    const baseUrl: string = "localhost:8000";

    let currentRoute: string;

    var connection: WebSocket = new WebSocket("ws://" + baseUrl + "/ws");

    onMount(() => {
        // conn = new WebSocket("ws://" + baseUrl + "/ws");

        if (currentRoute !== "/login" && currentRoute !== "/signup") {
        }
        console.log("This code runs on all routes");

        document.body.classList.add("js-enabled");
    });

    connection.onopen = function () {
        const authToken = localStorage.getItem("authToken") || "";
        connection.send(authToken);
        console.log("WebSocket connection established successfully.");
    };

    connection.onclose = function (evt) {
        console.log("Websocket connection closed");
        // var item = document.createElement("div");
        // item.innerHTML = "<b>Connection closed.</b>";
        // appendLog(item);
    };

    // Subscribe to the page store to know when the route changes
    $: {
        currentRoute = $page.url.pathname;
        console.log("Current route:", currentRoute);

        // You can add more code here that should run on every route change
    }
</script>

<slot />
