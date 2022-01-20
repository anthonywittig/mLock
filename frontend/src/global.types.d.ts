type AuditLogT = {
    id: string,
    entries: AuditLogEntriesT[],
}

type AuditLogEntriesT = {
    createdAt: string,
    log: string,
}

type DeviceLockCodeT = {
    code: string,
    mode: string,
    name: string,
    slot: number,
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
        lockCodes: [DeviceLockCode] | null,
        name: string,
        status: string,
    }
}

type DeviceManagedLockCodeT = {
    id: string,
    deviceId: string,
    code: string,
    note: string,
    reservation: DeviceManagedLockCodeReservationT,
    status: string,
    startAt: string,
    endAt: string,
}

type DeviceManagedLockCodeReservationT = {
    id: string,
    sync: boolean,
}

type UnitT = {
    id: string,
    name: string,
    propertyId: string,
    calendarUrl: string,
    updatedBy: string,
}
