export interface User {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  secondary_id: string;
  username: string;
  SentRequests: any | null;
  ReceivedRequests: any | null;
  Chats: Chat[] | null;
}

export interface Chat {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  secondary_id: string;
  name: string;
  is_private_chat: boolean;
  ConnectionRequests: any | null;
  users: User[];
}

export interface Message {
  dbId: string | null;
  localId: string;
  chatId: string;

  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;

  message: string;
  from: User;
}
