<script lang="ts">
    import { goto } from "$app/navigation";

    import { API_BASE_URL } from "$lib/config/api";

    let username = "";
    let password = "";
    let confirmPassword = "";
    let errorText = "";

    function showError(error: string) {
        errorText = error;
    }

    async function handleSubmit(event: Event) {
        event.preventDefault();

        showError("");

        if (password !== confirmPassword) {
            showError("Passwords didn't match");
            return;
        }

        const formData = {
            username,
            password,
            confirmPassword: confirmPassword,
        };

        try {
            const response = await fetch(`http://${API_BASE_URL}/signup`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(formData),
            });

            if (response.status === 201) {
                goto("/login");
            } else {
                const result = await response.json();
                showError(result.message || "An error occurred");
            }
        } catch (error) {
            console.error("Error:", error);
            showError("An unexpected error occurred");
        }
    }
</script>

<div class="columns is-centered">
    <div class="column is-one-third">
        <form
            on:submit={handleSubmit}
            class="has-background-dark box has-radius-medium"
        >
            <h2 class="title is-4 has-text-centered">Signup</h2>

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

            <div class="field">
                <label class="label" for="confirm-password"
                    >Confirm Password:</label
                >
                <div class="control">
                    <input
                        class="input"
                        type="password"
                        id="confirm-password"
                        bind:value={confirmPassword}
                        required
                    />
                </div>
            </div>

            <div class="field is-grouped is-grouped-centered">
                <div class="control">
                    <button class="button is-primary" type="submit"
                        >Sign Up</button
                    >
                </div>
            </div>

            {#if errorText}
                <span id="errorText" class="has-text-danger">{errorText}</span>
            {/if}
            <a href="/login" class="has-text-centered">Login</a>
        </form>
    </div>
</div>
