export interface Maintenance {
    id: string,
    equipmentId: string,
    equipmentCode: string,
    equipmentType: string,
    hmMaintenance: number,
    maintenanceDate: string,
    description?: string,
}

export interface MaintenanceRequest {
    equipmentId: string,
    hmMaintenance: number,
    maintenanceDate: string,
    description?: string,
}