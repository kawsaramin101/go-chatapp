<script lang="ts">
    import { API_BASE_URL } from "$lib/config/api";

    let username: string = "";
    let password: string = "";
    let error: string = "";

    const baseUrl: string = "http://localhost:8000";

    async function handleSubmit(event: Event): Promise<void> {
        event.preventDefault();
        error = ""; // Clear any previous error

        try {
            const response: Response = await fetch(
                `http://${API_BASE_URL}/login`,
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({ username, password }),
                },
            );

            if (response.status === 200) {
                const result: {
                    token: string;
                    username: string;
                    [key: string]: any;
                } = await response.json();
                localStorage.setItem("authToken", result.token);
                localStorage.setItem("username", result.username);
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

<div class="columns is-centered">
    <div class="column is-one-third">
        <form
            on:submit={handleSubmit}
            class="has-background-dark is-shadowless box has-radius-medium"
        >
            <h2 class="title is-4 has-text-centered">Login</h2>

            <div class="field">
                <label class="label" for="username">Username:</label>
                <div class="control">
                    <input
                        class="input"
                        type="text"
                        id="username"
                        bind:value={username}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <label class="label" for="password">Password:</label>
                <div class="control">
                    <input
                        class="input"
                        type="password"
                        id="password"
                        bind:value={password}
                        required
                    />
                </div>
            </div>

            <div class="field is-grouped is-grouped-centered">
                <div class="control">
                    <button class="button is-primary" type="submit"
                        >Login</button
                    >
                </div>
            </div>

            {#if error}
                <span id="errorText" class="has-text-danger">{error}</span>
            {/if}

            <a href="/signup">Signup</a>
        </form>
    </div>
</div>
