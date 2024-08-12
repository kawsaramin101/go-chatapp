CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    secondary_id TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    connection_request_from INTEGER,
    connection_request_to INTEGER,
    FOREIGN KEY (connection_request_from) REFERENCES user(id),
    FOREIGN KEY (connection_request_to) REFERENCES user(id)
);


CREATE TABLE IF NOT EXISTS connection (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    secondary_id TEXT UNIQUE NOT NULL,
    send_by INTEGER NOT NULL,
    send_to INTEGER NOT NULL,
    is_accepted BOOLEAN,
    FOREIGN KEY (send_by) REFERENCES user(id),
    FOREIGN KEY (send_to) REFERENCES user(id),
    UNIQUE(send_by, send_to)
);
