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

<form on:submit={handleSubmit}>
    <h2>Signup</h2>
    <label for="username">Username:</label>
    <input type="text" id="username" bind:value={username} required />

    <label for="password">Password:</label>
    <input type="password" id="password" bind:value={password} required />

    <label for="confirm-password">Confirm Password:</label>
    <input
        type="password"
        id="confirm-password"
        bind:value={confirmPassword}
        required
    />

    <button type="submit">Sign Up</button>

    <span>{errorText}</span>
</form>

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
        align-items: center;
    }
    h2 {
        align-self: center;
    }
    label,
    button {
        margin-bottom: 10px;
    }
</style>
