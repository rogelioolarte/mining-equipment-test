export interface Equipment {
    id: string,
    code: string,
    type: string,
    currentHourMeter: number,
    maintenanceIntervalHours: number,
    hoursSinceMaintenance: number,
    lastMaintenanceHourMeter: number,
    maintenanceStatus: string
}

export interface EquipmentRequest {
    code: string,
    typeId: string,
    current_hour_meter: number,
}