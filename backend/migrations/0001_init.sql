CREATE SCHEMA app;

CREATE TABLE IF NOT EXISTS app.equipment_type (
    id uuid PRIMARY KEY NOT NULL,
    name varchar(100) NOT NULL,
    maintenance_hours integer NOT NULL
);

CREATE TABLE IF NOT EXISTS app.equipment (
    id uuid PRIMARY KEY NOT NULL,
    type_id uuid NOT NULL REFERENCES app.equipment_type (id),
    code VARCHAR(50) NOT NULL,
    current_h_meter integer NOT NULL,
    last_mh_meter integer NOT NULL
);

CREATE TABLE IF NOT EXISTS app.maintenance (
    id uuid PRIMARY KEY NOT NULL,
    equipment_id uuid NOT NULL REFERENCES app.equipment (id),
    hm_maintenance integer NOT NULL,
    maintenance_date timestamp NOT NULL,
    description varchar(100)
);