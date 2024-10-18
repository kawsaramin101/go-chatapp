<script lang="ts">
    import { connectionRequestStore } from "$lib/stores/connectionRequests";
    console.log($connectionRequestStore);
</script>

<div class="columns is-centered">
    <div class="column is-half">
        <div class="container has-text-centered">
            <div class="tabs is-toggle is-fullwidth">
                <ul>
                    <li><a href="/">Inbox</a></li>
                    <li class="is-active">
                        <a href="message_requests">Requests</a>
                    </li>
                </ul>
            </div>
            <ol>
                {#each $connectionRequestStore as connectionRequest}
                    <li class="mt-2">
                        <b> {connectionRequest.sendBy["username"]}</b> invited
                        you to Join chat
                        <a href="/chat/{connectionRequest.chat.ID}">
                            {connectionRequest.chat.ID}
                        </a>
                    </li>
                {/each}
            </ol>
        </div>
    </div>
</div>
