<script lang="ts">
    import { wsStore } from "$lib/stores/ws";

    function addUser(event: SubmitEvent) {
        event.preventDefault();

        const formData = new FormData(event.target as HTMLFormElement);

        const sendingData = {
            action: "CREATECHAT",
            data: {
                username: formData.get("username"),
            },
        };
        $wsStore!.send(JSON.stringify(sendingData));
    }
</script>

<form on:submit={addUser}>
    <input type="text" name="username" placeholder="Username" />
    <button type="submit">GO</button>
</form>

<h3>Create Group Chat</h3>
<form></form>
