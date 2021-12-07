type AuditLogT = {
    id: string,
    entries: AuditLogEntriesT[],
}

type AuditLogEntriesT = {
    createdAt: string,
    log: string,
}

type DeviceT = {
    id: string,
    propertyId: string,
    unitId: string | null,
    lastRefreshedAt: string,
    lastWentOfflineAt: string | null,
    lastWentOnlineAt: string | null,
    managedLockCodes: DeviceManagedLockCode[],
    rawDevice: {
        battery: {
            batteryPowered: boolean,
            level: number,
        },
        categoryId: string,
        lockCodes: [
            {
                code: string,
                mode: string,
                name: string,
                slot: number,
            }
        ] | null,
        name: string,
        status: string,
    }
}

type DeviceManagedLockCode = {
    id: string,
    deviceId: string,
    code: string,
    note: string,
    status: string,
    startAt: string,
    endAt: string,
}

type UnitT = {
    id: string,
    name: string,
    propertyId: string,
    calendarUrl: string,
    updatedBy: string,
}