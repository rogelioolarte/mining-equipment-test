CREATE SCHEMA app;

CREATE TABLE IF NOT EXISTS app.equipment_type (
    id uuid PRIMARY KEY NOT NULL,
    name varchar(100) UNIQUE NOT NULL,
    maintenance_hours integer NOT NULL
);

CREATE TABLE IF NOT EXISTS app.equipment (
    id uuid PRIMARY KEY NOT NULL,
    type_id uuid NOT NULL REFERENCES app.equipment_type (id) ON DELETE CASCADE,
    code VARCHAR(50) UNIQUE NOT NULL,
    current_h_meter integer NOT NULL
);

CREATE TABLE IF NOT EXISTS app.maintenance (
    id uuid PRIMARY KEY NOT NULL,
    equipment_id uuid NOT NULL REFERENCES app.equipment (id) ON DELETE CASCADE,
    hm_maintenance integer NOT NULL,
    maintenance_date timestamp NOT NULL,
    description varchar(100)
);

CREATE INDEX IF NOT EXISTS idx_maintenance_equipment_date
ON app.maintenance (equipment_id, maintenance_date DESC, id DESC);