// src/stores/websocket.ts
import { writable, get } from "svelte/store";
import type { Writable } from "svelte/store";
import { page } from "$app/stores";

import Toastify from "toastify-js";
import "toastify-js/src/toastify.css";

import { chats } from "$lib/stores/chats";
import { API_BASE_URL } from "$lib/config/api";

export const wsStore: Writable<WebSocket | null> = writable(null);
let wsInstance: WebSocket | null = null;

export function initializeWebSocket(): WebSocket {
    if (wsInstance) return wsInstance;

    wsInstance = new WebSocket("ws://" + API_BASE_URL + "/ws");

    wsInstance.onopen = () => {
        console.log("WebSocket connection established successfully.");
        wsStore.set(wsInstance);

        const authToken = localStorage.getItem("authToken") || "";
        wsInstance!.send(authToken);
    };

    wsInstance.onmessage = (event: MessageEvent) => {
        const data = JSON.parse(event.data);
        switch (data["action"]) {
            case "ERROR_USER_NOT_FOUND":
            case "ERROR_SERVER_ERROR":
            case "ERROR_INVALID_PAYLOAD":
                Toastify({
                    text: data["message"],
                    duration: 3000,
                    close: true,
                    gravity: "top", // `top` or `bottom`
                    position: "center", // `left`, `center` or `right`
                    stopOnFocus: true, // Prevents dismissing of toast on hover
                    onClick: function () {}, // Callback after click
                }).showToast();
                break;

            case "INITIAL_DATA":
                chats.setChats(data["data"]["chats"]);
                break;

            default:
                // Handle any other actions if needed
                break;
        }
    };

    wsInstance.onerror = (error: Event) => {
        console.error("WebSocket error:", error);
    };

    wsInstance.onclose = (event: CloseEvent) => {
        console.log("WebSocket connection closed", event);
        wsStore.set(null);
        wsInstance = null;

        const currentRoute = get(page).url.pathname;
        if (
            !event.wasClean &&
            currentRoute !== "/login" &&
            currentRoute !== "/signup"
        ) {
            setTimeout(() => {
                wsInstance = initializeWebSocket();
                wsStore.set(wsInstance);
            }, 2000);
        }
    };

    return wsInstance;
}

export function closeWebSocket(): void {
    if (wsInstance && wsInstance.readyState === WebSocket.OPEN) {
        wsInstance.close();
    }
    wsStore.set(null);
    wsInstance = null;
}
