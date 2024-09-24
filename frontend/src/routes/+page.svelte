<script>
    import { getContext } from "svelte";
    import { chats } from "$lib/stores/chats";

    // Get the function from +layout.svelte
    const addUser = getContext("addUser");
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
