import { writable } from "svelte/store";
import type { Message } from "$lib/models";

export const messageStore = writable<Message[]>([]);

export const messageStoreOperations = {
    subscribe: messageStore.subscribe,
    set: messageStore.set,
    update: messageStore.update,

    addMessage: (newMessage: Message) => {
        messageStore.update((messages) => [...messages, newMessage]);
    },

    // removeMessage: (id: string) => {

    //   messageStore.update((messages) =>
    //     messages.filter((message) => message.id? !== id),
    //   );
    // },

    // updateMessage: (updatedMessage: Message) => {
    //   messageStore.update((messages) =>
    //     messages.map((message) =>
    //       message.localId === updatedMessage.localId ? updatedMessage : message,
    //     ),
    //   );
    // },
    setMessages: (newMessages: Message[]) => {
        messageStore.set(newMessages);
    },
};
