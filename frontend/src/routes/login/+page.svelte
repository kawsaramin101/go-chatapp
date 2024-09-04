<script>
    let username = "";
    let password = "";
    let error = "";

    async function handleSubmit(event) {
        event.preventDefault();
        error = ""; // Clear any previous error

        try {
            const response = await fetch("/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded",
                },
                body: new URLSearchParams({ username, password }),
            });

            if (response.status === 200) {
                window.location.href = "/";
            } else {
                const result = await response.json();
                error = result.message || "An error occurred";
            }
        } catch (err) {
            console.error("Error:", err);
            error = "An error occurred";
        }
    }
</script>

<main>
    <form on:submit={handleSubmit}>
        <h2>Login</h2>
        <label for="username">Username:</label>
        <input type="text" id="username" bind:value={username} required />

        <label for="password">Password:</label>
        <input type="password" id="password" bind:value={password} required />

        <button type="submit">Login</button>

        {#if error}
            <span id="errorText">{error}</span>
        {/if}
    </form>
</main>

<style>
    :global(html),
    :global(body) {
        height: 100%;
        margin: 0;
        display: flex;
        justify-content: center;
        align-items: center;
    }
    form {
        display: flex;
        flex-direction: column;
    }
    h2 {
        align-self: center;
    }
    input,
    button {
        margin-bottom: 10px;
    }
</style>
