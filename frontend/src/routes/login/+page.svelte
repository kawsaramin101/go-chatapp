<script lang="ts">
    let username: string = "";
    let password: string = "";
    let error: string = "";

    const baseUrl: string = "http://localhost:8000";

    async function handleSubmit(event: Event): Promise<void> {
        event.preventDefault();
        error = ""; // Clear any previous error

        try {
            const response: Response = await fetch(`${baseUrl}/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, password }),
            });

            if (response.status === 200) {
                const result: { token: string; [key: string]: any } =
                    await response.json();
                const token: string = result.token;
                localStorage.setItem("authToken", token);
                window.location.href = "/";
            } else {
                const result: { message?: string; [key: string]: any } =
                    await response.json();
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
