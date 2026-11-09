
-- ============================================================
-- 0002_data.sql
-- Default data for Mining Equipment Registration Application
-- UUID version: 7
-- ============================================================

-- ------------------------------------------------------------
-- 1. EQUIPMENT TYPES
-- ------------------------------------------------------------

INSERT INTO app.equipment_type (
    id,
    name,
    maintenance_hours
)
VALUES
    (
        '01a1b874-6c00-7def-b20f-edcba9874421',
        'Camión',
        250
    ),
    (
        '01a1b874-6c01-7e00-b20f-edcba9874440',
        'Excavadora',
        200
    ),
    (
        '01a1b874-6c02-7e11-b20f-edcba987445f',
        'Perforadora',
        150
    );

-- ------------------------------------------------------------
-- 2. EQUIPMENT
-- ------------------------------------------------------------

INSERT INTO app.equipment (
    id,
    type_id,
    code,
    current_h_meter
)
VALUES
    -- Trucks
    (
        '01a1b874-6c03-7e22-b20f-edcba987447e',
        '01a1b874-6c00-7def-b20f-edcba9874421',
        'CAM-001',
        1250
    ),
    (
        '01a1b874-6c04-7e33-b20f-edcba987449d',
        '01a1b874-6c00-7def-b20f-edcba9874421',
        'CAM-002',
        980
    ),
    (
        '01a1b874-6c05-7e44-b20f-edcba98744bc',
        '01a1b874-6c00-7def-b20f-edcba9874421',
        'CAM-003',
        420
    ),

    -- Excavators
    (
        '01a1b874-6c06-7e55-b20f-edcba98744db',
        '01a1b874-6c01-7e00-b20f-edcba9874440',
        'EXC-001',
        860
    ),
    (
        '01a1b874-6c07-7e66-b20f-edcba98744fa',
        '01a1b874-6c01-7e00-b20f-edcba9874440',
        'EXC-002',
        1240
    ),
    (
        '01a1b874-6c08-7e77-b20f-edcba9874519',
        '01a1b874-6c01-7e00-b20f-edcba9874440',
        'EXC-003',
        600
    ),

    -- Drilling machines
    (
        '01a1b874-6c09-7e88-b20f-edcba9874538',
        '01a1b874-6c02-7e11-b20f-edcba987445f',
        'PER-001',
        540
    ),
    (
        '01a1b874-6c0a-7e99-b20f-edcba9874557',
        '01a1b874-6c02-7e11-b20f-edcba987445f',
        'PER-002',
        310
    ),
    (
        '01a1b874-6c0b-7eaa-b20f-edcba9874576',
        '01a1b874-6c02-7e11-b20f-edcba987445f',
        'PER-003',
        180
    );

-- ------------------------------------------------------------
-- 3. MAINTENANCE HISTORY
-- ------------------------------------------------------------

INSERT INTO app.maintenance (
    id,
    equipment_id,
    hm_maintenance,
    maintenance_date,
    description
)
VALUES
    -- CAM-001: Multiple maintenance records.
    -- Last maintenance: 1100. Current: 1250. Difference: 150.
    (
        '01a1b874-6c0c-7ebb-b20f-edcba9874595',
        '01a1b874-6c03-7e22-b20f-edcba987447e',
        600,
        '2026-01-15 08:30:00Z',
        'Mantenimiento preventivo general'
    ),
    (
        '01a1b874-6c0d-7ecc-b20f-edcba98745b4',
        '01a1b874-6c03-7e22-b20f-edcba987447e',
        850,
        '2026-05-20 10:00:00Z',
        'Cambio de filtros y lubricación'
    ),
    (
        '01a1b874-6c0e-7edd-b20f-edcba98745d3',
        '01a1b874-6c03-7e22-b20f-edcba987447e',
        1100,
        '2026-09-12 09:00:00Z',
        'Mantenimiento preventivo de motor'
    ),

    -- CAM-002: Requires maintenance.
    -- Last maintenance: 700. Current: 980. Difference: 280.
    (
        '01a1b874-6c0f-7eee-b20f-edcba98745f2',
        '01a1b874-6c04-7e33-b20f-edcba987449d',
        700,
        '2026-08-10 08:00:00Z',
        'Mantenimiento preventivo general'
    ),

    -- EXC-001: Operational.
    -- Last maintenance: 700. Current: 860. Difference: 160.
    (
        '01a1b874-6c10-7eff-b20f-edcba9874611',
        '01a1b874-6c06-7e55-b20f-edcba98744db',
        400,
        '2026-02-10 09:00:00Z',
        'Inspección y mantenimiento de componentes'
    ),
    (
        '01a1b874-6c11-7f10-b20f-edcba9874630',
        '01a1b874-6c06-7e55-b20f-edcba98744db',
        700,
        '2026-08-18 11:00:00Z',
        'Cambio de aceite y filtros'
    ),

    -- EXC-002: Requires maintenance.
    -- Last maintenance: 1000. Current: 1240. Difference: 240.
    (
        '01a1b874-6c12-7f21-b20f-edcba987464f',
        '01a1b874-6c07-7e66-b20f-edcba98744fa',
        1000,
        '2026-07-25 08:30:00Z',
        'Mantenimiento preventivo general'
    ),

    -- EXC-003: Recently maintained.
    -- Last maintenance: 600. Current: 600. Difference: 0.
    (
        '01a1b874-6c13-7f32-b20f-edcba987466e',
        '01a1b874-6c08-7e77-b20f-edcba9874519',
        600,
        '2026-10-15 10:00:00Z',
        'Mantenimiento preventivo completo'
    ),

    -- PER-001: Requires maintenance.
    -- Last maintenance: 390. Current: 540. Difference: 150.
    (
        '01a1b874-6c14-7f43-b20f-edcba987468d',
        '01a1b874-6c09-7e88-b20f-edcba9874538',
        240,
        '2026-03-12 08:00:00Z',
        'Inspección del sistema de perforación'
    ),
    (
        '01a1b874-6c15-7f54-b20f-edcba98746ac',
        '01a1b874-6c09-7e88-b20f-edcba9874538',
        390,
        '2026-08-22 09:30:00Z',
        'Cambio de componentes y lubricación'
    ),

    -- PER-002: Operational.
    -- Last maintenance: 200. Current: 310. Difference: 110.
    (
        '01a1b874-6c16-7f65-b20f-edcba98746cb',
        '01a1b874-6c0a-7e99-b20f-edcba9874557',
        200,
        '2026-09-05 08:30:00Z',
        'Mantenimiento preventivo general'
    );
