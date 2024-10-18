import { writable } from "svelte/store";
import type { ConnectionRequest } from "$lib/models";

export const connectionRequestStore = writable<ConnectionRequest[]>([]);

export const chats = {
    subscribe: connectionRequestStore.subscribe,
    set: connectionRequestStore.set,
    update: connectionRequestStore.update,

    addChatRequest: (newConnectionRequest: ConnectionRequest) => {
        connectionRequestStore.update((connectionRequests) => [
            ...connectionRequests,
            newConnectionRequest,
        ]);
    },

    // removeChatRequest: (id: string) => {
    //   chatRequestStore.update((chatRequests) => chatRequests.filter((chatRequest) => chatRequest.secondaryID !== id));
    // },

    // updateChat: (updatedChat: ChatRequest) => {
    //   chatRequestStore.update((chatRequests) =>
    //     chatRequests.map((chatRequest) => (chatRequest.secondaryID === updatedChat.secondaryID ? updatedChat : chatRequest)),
    //   );
    // },
    // setChats: (newChats: Chat[]) => {
    //   chatRequestStore.set(newChats);
    // },
};
