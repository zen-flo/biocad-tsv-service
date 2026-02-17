-- Migration: create messages table
-- Stores all messages parsed from TSV files

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE "messages" (
                            id uuid PRIMARY KEY DEFAULT gen_random_uuid(),          -- unique identifier
                            mqtt text,                                              -- optional MQTT broker or topic
                            unit_guid uuid NOT NULL,                                -- device GUID
                            msg_id text,                                            -- message ID from TSV
                            text text,                                              -- message text
                            context text,                                           -- environment or context
                            class text,                                             -- message class: alarm, warning, info, event, command
                            level int,                                              -- message level (integer)
                            area text,                                              -- variable area (HR, IR, I, C)
                            addr text,                                              -- variable address in controller
                            block text NULL,                                             -- use as block start
                            type text,                                              -- type
                            bit text NULL,                                               -- bit number in register
                            invert_bit text NULL,                                        -- inverted bit flag
                            created_at timestamp NOT NULL DEFAULT now()             -- timestamp of creation
);

CREATE INDEX idx_messages_unit_guid ON "messages"(unit_guid);
