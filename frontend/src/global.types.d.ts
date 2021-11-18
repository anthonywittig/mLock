type DeviceT = {
    id: string,
    propertyId: string,
    unitId: string | null,
    lastRefreshedAt: string,
    lastWentOfflineAt: string | null,
    lastWentOnlineAt: string | null,
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
    startAt: Date,
    endAt: Date,
}

type UnitT = {
    id: string,
    name: string,
    propertyId: string,
    calendarUrl: string,
    updatedBy: string,
}