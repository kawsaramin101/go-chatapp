// src/stores/websocket.ts
import { writable } from "svelte/store";
import type { Writable } from "svelte/store";
import { chats } from "$lib/stores/chats";

interface WebSocketStore {
    subscribe: Writable<WebSocket | null>["subscribe"];
    get: () => WebSocket; // Return null if WebSocket is not initialized
    close: () => void;
    set: () => void;
}

let websocketInstance: WebSocket | null = null;

const createWebSocket = (): WebSocket => {
    const baseUrl: string = "localhost:8000";
    const ws = new WebSocket("ws://" + baseUrl + "/ws");

    ws.onopen = function () {
        console.log("WebSocket connection established successfully.");
        const authToken = localStorage.getItem("authToken") || "";
        ws.send(authToken);
    };

    ws.onmessage = (event: MessageEvent) => {
        console.log(event.data);
        const data = JSON.parse(event.data);
        switch (data["action"]) {
            case "ERROR_USER_NOT_FOUND":
            case "ERROR_SERVER_ERROR":
            case "ERROR_INVALID_PAYLOAD":
                alert(data["message"]);
                break;

            case "CHAT_CREATED":
                alert("Chat created");
                setTimeout(() => {}, 3000);
                break;

            case "INITIAL_DATA":
                chats.setChats(data["data"]["chats"]);
                break;

            case "MESSAGE":

            default:
                // Handle any other actions if needed
                break;
        }
    };

    ws.onclose = function (event) {
        console.log("WebSocket connection closed", event);
        // Retry logic on unclean closure
        if (!event.wasClean) {
            setTimeout(() => {
                websocketInstance = createWebSocket(); // Create a new WebSocket instance
                websocketStore.set(websocketInstance); // Update store with the new instance
            }, 2000);
        }
    };

    return ws;
};

const websocketStore = (() => {
    const { subscribe, set }: Writable<WebSocket | null> = writable(null);

    const get = (): WebSocket => {
        if (
            websocketInstance === null ||
            websocketInstance.readyState === WebSocket.CLOSED
        ) {
            console.log("run");
            websocketInstance = createWebSocket(); // Create a new WebSocket if none exists
            set(websocketInstance); // Update the store
        }
        return websocketInstance; // Return the current WebSocket instance
    };

    const close = (): void => {
        if (
            websocketInstance &&
            websocketInstance.readyState === WebSocket.OPEN
        ) {
            websocketInstance.close(); // Close the WebSocket
            websocketInstance = null; // Clear the instance
            set(null); // Clear the store
        }
    };

    return {
        subscribe,
        get,
        close,
        set,
    };
})();

export const websocket: WebSocketStore = websocketStore;
