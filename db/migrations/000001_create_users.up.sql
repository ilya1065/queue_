CREATE TABLE users (
                       id   INTEGER PRIMARY KEY,
                       name TEXT NOT NULL
);

CREATE TABLE schedule_items (
                                id          INTEGER PRIMARY KEY,
                                name TEXT not null,
                                description TEXT,
                                start_date  timestamp NOT NULL,
                                end_date    timestamp NOT NULL,
                                external_id TEXT NOT NULL,
                                UNIQUE (external_id,start_date)
);

CREATE TABLE records (
                         id               INTEGER PRIMARY KEY,
                         user_id          INTEGER NOT NULL,
                         schedule_item_id INTEGER NOT NULL,
                         createdAt TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
                        /* body             TEXT,*/

                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                         FOREIGN KEY (schedule_item_id) REFERENCES schedule_items(id) ON DELETE CASCADE
                         UNIQUE (user_id,schedule_item_id)
);

CREATE INDEX idx_records_user_id ON records(user_id);
CREATE INDEX idx_records_schedule_item_id ON records(schedule_item_id);

