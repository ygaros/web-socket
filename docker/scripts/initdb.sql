create table IF NOT EXISTS client (
    id TEXT PRIMARY KEY,
    name text not null unique
);

create table IF NOT EXISTS notification (
    id text PRIMARY KEY,
    message text not null,
    creation_date text not null,
    sender text not null
);

create table IF NOT EXISTS client_notifications (
    client_id text,
    notification_id text,
    notification_status text,
    PRIMARY KEY (client_id, notification_id),
    FOREIGN KEY (client_id)
        REFERENCES client (id)
            ON DELETE CASCADE
            ON UPDATE NO ACTION,
    FOREIGN KEY (notification_id)
        REFERENCES notification (id)
            ON DELETE CASCADE
            ON UPDATE NO ACTION
);
