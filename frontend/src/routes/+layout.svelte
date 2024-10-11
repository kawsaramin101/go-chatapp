<script lang="ts">
    import { page } from "$app/stores";
    import { beforeNavigate, goto } from "$app/navigation";
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
