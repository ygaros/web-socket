#!/usr/bin/env sh

sqlite3 -batch "$PWD/docker/scripts/websocket_database.sqlite" <"$PWD/docker/scripts/initdb.sql"