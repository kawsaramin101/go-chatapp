import { writable } from "svelte/store";
import type { Chat } from "$lib/models";

const chatsStore = writable<Chat[]>([]);

export const chats = {
  subscribe: chatsStore.subscribe,
  set: chatsStore.set,
  update: chatsStore.update,

  addChat: (newChat: Chat) => {
    chatsStore.update((chats) => [...chats, newChat]);
  },

  removeChat: (id: number) => {
    chatsStore.update((chats) => chats.filter((chat) => chat.ID !== id));
  },

  updateChat: (updatedChat: Chat) => {
    chatsStore.update((chats) =>
      chats.map((chat) => (chat.ID === updatedChat.ID ? updatedChat : chat)),
    );
  },
  setChats: (newChats: Chat[]) => {
    chatsStore.set(newChats);
  },
};
