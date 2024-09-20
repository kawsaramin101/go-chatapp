<script>
    import { getContext } from "svelte";

    // Get the function from +layout.svelte
    const addUser = getContext("addUser");
</script>

<form on:submit={addUser}>
    <input type="text" name="username" placeholder="Username" />
    <button type="submit">GO</button>
</form>
<a href="/login">Login</a>
