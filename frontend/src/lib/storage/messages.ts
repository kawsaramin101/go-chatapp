import type { Message } from "$lib/models";

let db: IDBDatabase | null = null;

const DB_NAME = "messageDB";
const DB_VERSION = 4;

export function openDatabase(): Promise<IDBDatabase> {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, DB_VERSION);

        request.onupgradeneeded = function (event) {
            const db = (event.target as IDBOpenDBRequest).result;
            if (!db.objectStoreNames.contains("messages")) {
                const messagesStore = db.createObjectStore("messages", {
                    keyPath: "id",
                    autoIncrement: true,
                });
                // messagesStore.createIndex("dbId", "dbId", {
                //     unique: true,
                // });
                messagesStore.createIndex("dbSecondaryId", "dbSecondaryId", {
                    unique: false,
                });
                messagesStore.createIndex("chatId", "chatId", {
                    unique: false,
                });
                messagesStore.createIndex("from", "from", {
                    unique: false,
                });
                messagesStore.createIndex("content", "content", {
                    unique: false,
                });
                messagesStore.createIndex("createdAt", "createdAt", {
                    unique: false,
                });
                messagesStore.createIndex(
                    "chatId_createdAt",
                    ["chatId", "createdAt"],
                    {
                        unique: false,
                    },
                );
            }
        };

        request.onsuccess = function (event) {
            db = (event.target as IDBOpenDBRequest).result;
            resolve(db);
        };

        request.onerror = function (event) {
            reject((event.target as IDBRequest).error);
        };
    });
}

export function getDatabase(): Promise<IDBDatabase> {
    if (db) {
        return Promise.resolve(db);
    } else {
        return openDatabase();
    }
}

export async function getAllMessagesFromAChat(
    chatId: number,
): Promise<Message[]> {
    try {
        const db = await getDatabase();
        const transaction = db.transaction(["messages"], "readonly");
        const store = transaction.objectStore("messages");
        const index = store.index("chatId_createdAt");
        console.log(chatId);
        const range = IDBKeyRange.bound(
            [chatId, new Date(0)],
            [chatId, new Date()],
        );

        return new Promise((resolve, reject) => {
            const request = index.openCursor(range, "prev");
            const messages: Message[] = [];

            request.onsuccess = function (event) {
                const cursor = (event.target as IDBRequest<IDBCursorWithValue>)
                    .result;
                if (cursor) {
                    messages.push(cursor.value);
                    cursor.continue();
                } else {
                    resolve(messages);
                }
            };

            request.onerror = function (event) {
                reject((event.target as IDBRequest).error);
            };
        });
    } catch (error) {
        console.error("Failed to retrieve messages:", error);
        throw error;
    }
}

export async function addMessageToStore(message: Message): Promise<Message> {
    try {
        const db = await getDatabase();
        const transaction = db.transaction("messages", "readwrite");
        const store = transaction.objectStore("messages");

        console.log(message);

        const request: IDBRequest<IDBValidKey> = store.add(message);

        return new Promise((resolve, reject) => {
            request.onsuccess = function (event) {
                const id = (event.target as IDBRequest<IDBValidKey>)
                    .result as number;
                console.log(`Message added with id: ${id}`);

                // You can resolve with the original message + ID if needed
                resolve({ ...message, id });
            };

            request.onerror = function (event) {
                console.error((event.target as IDBRequest).error, event.target);
                reject((event.target as IDBRequest).error);
            };
        });
    } catch (error) {
        console.error("Failed to add message:", error);
        throw error;
    }
}

// export async function getNote(id, secondaryId) {
//     if (id === undefined && secondaryId === undefined) {
//         throw new Error("ID or secondaryId must be provided to get a note");
//     }

//     const db = await getDatabase();
//     const transaction = db.transaction(["notes"], "readwrite");
//     const store = transaction.objectStore("notes");

//     try {
//         let getRequest;
//         if (id) {
//             getRequest = store.get(id);
//         } else {
//             const index = store.index("secondaryId");
//             getRequest = index.get(secondaryId);
//         }

//         return new Promise((resolve, reject) => {
//             getRequest.onsuccess = function (event) {
//                 if (request.result) {
//                     resolve(request.result);
//                 } else {
//                     reject(new Error("Note not found"));
//                 }
//             };

//             getRequest.onerror = function (event) {
//                 reject(event.target.error);
//             };
//         });
//     } catch (err) {
//         reject(err);
//     }
// }

// export async function addNote(note) {
//     try {
//         const db = await getDatabase();
//         const transaction = db.transaction(["notes"], "readwrite");
//         const store = transaction.objectStore("notes");

//         const lastAddedNoteOrder = parseInt(
//             localStorage.getItem("lastAddedNoteOrder") ?? "0",
//             10,
//         );

//         const noteData = {
//             secondaryId: note.hasOwnProperty("secondaryId")
//                 ? note.secondaryId
//                 : uniqid(),
//             title: note.title,
//             content: note.content,
//             pinned: note.hasOwnProperty("pinned") ? note.pinned : false,
//             order: lastAddedNoteOrder + 1,
//             createdAt: note.hasOwnProperty("createdAt")
//                 ? note.createdAt
//                 : new Date().toISOString(),
//             updatedAt: note.hasOwnProperty("updatedAt")
//                 ? note.updatedAt
//                 : new Date().toISOString(),
//         };

//         const request = store.add(noteData);

//         localStorage.setItem("lastAddedNoteOrder", lastAddedNoteOrder + 1);

//         return new Promise((resolve, reject) => {
//             request.onsuccess = function (event) {
//                 const id = event.target.result;
//                 const getRequest = store.get(id);

//                 getRequest.onsuccess = function () {
//                     resolve(getRequest.result);
//                 };

//                 getRequest.onerror = function (event) {
//                     reject(event.target.error);
//                 };
//             };

//             request.onerror = function (event) {
//                 if (event.target.error.name === "ConstraintError") {
//                     console.error(
//                         "ConstraintError: A note with the same unique field value already exists.",
//                     );
//                     fixSecondaryIdError();

//                     const note = getNote(undefined, note.secondaryId);
//                     updateNote({ id: note.id, ...note });
//                 }
//                 reject(event.target.error);
//             };
//         });
//     } catch (error) {
//         console.error("Failed to add note:", error);
//         throw error;
//     }
// }

// async function fixSecondaryIdError(secondaryId) {
//     const db = await getDatabase();
//     const transaction = db.transaction(["notes"], "readwrite");
//     const store = transaction.objectStore("notes");

//     try {
//         const note = getNote(undefined, secondaryId);
//         updateNote({});
//     } catch (err) {
//         reject(err);
//     }
// }

// export async function deleteNote(id) {
//     try {
//         const db = await getDatabase();
//         const transaction = db.transaction(["notes"], "readwrite");
//         const objectStore = transaction.objectStore("notes");

//         const deleteRequest = objectStore.delete(id);

//         deleteRequest.onsuccess = function (event) {
//             return true;
//         };

//         deleteRequest.onerror = function (event) {
//             console.error("Error deleting entry:", event.target.error);
//         };
//         transaction.onerror = function (event) {
//             console.error("Error:", event.target.error);
//         };
//     } catch (error) {
//         console.error("Error:", error);
//     }
// }

// export async function updateNote(note) {
//     try {
//         const db = await getDatabase();
//         const transaction = db.transaction(["notes"], "readwrite");
//         const objectStore = transaction.objectStore("notes");

//         const getRequest = objectStore.get(note.id);

//         getRequest.onsuccess = function (event) {
//             const entry = event.target.result;

//             if (note.title) entry.title = note.title;
//             if (note.content) entry.content = note.content;
//             if (note.hasOwnProperty("pinned")) entry.pinned = note.pinned;
//             if (note.order) entry.order = note.order;

//             if (note.title || note.content)
//                 entry.undatedAt = new Date().toISOString();

//             const putRequest = objectStore.put(entry);
//             putRequest.onerror = function (event) {
//                 console.error("Error updating entry:", event.target.error);
//             };
//         };

//         getRequest.onerror = function (event) {
//             console.error("Error retrieving entry:", event.target.error);
//         };

//         transaction.onerror = function (event) {
//             console.error("Error in transaction:", event.target.error);
//         };
//     } catch (error) {
//         console.error("Error:", error);
//     }
// }

export default { getAllMessagesFromAChat, addMessageToStore };
