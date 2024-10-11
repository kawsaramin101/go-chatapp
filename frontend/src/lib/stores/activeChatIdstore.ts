import { writable } from "svelte/store";
import type { Writable } from "svelte/store";

export const activeChatId: Writable<number | null> = writable(null);
