// src/stores/websocket.ts
import { writable } from "svelte/store";
import type { Writable } from "svelte/store";

interface WebSocketStore {
    subscribe: Writable<WebSocket | null>["subscribe"];
    get: () => WebSocket | null; // Return null if WebSocket is not initialized
    close: () => void;
}

let websocketInstance: WebSocket | null = null;

const createWebSocket = (): WebSocket => {
    const ws = new WebSocket("ws://your-websocket-url");

    ws.onopen = function () {
        console.log("WebSocket connection established successfully.");
        const authToken = localStorage.getItem("authToken") || "";
        ws.send(authToken);
    };

    ws.onmessage = (event: MessageEvent) => {
        console.log("Message received:", event.data);
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

    const get = (): WebSocket | null => {
        if (
            websocketInstance === null ||
            websocketInstance.readyState === WebSocket.CLOSED
        ) {
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
