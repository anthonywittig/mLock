type DeviceT = {
    id: string,
    propertyId: string,
    unitId: string | null,
    battery: {
        lastUpdatedAt: string | null,
        level: string,
    },
    lastRefreshedAt: string,
    lastWentOfflineAt: string | null,
    rawDevice: {
        categoryId: string,
        name: string,
        status: string,
    }
}

type UnitT = {
    id: string,
    name: string,
    propertyId: string,
    calendarUrl: string,
    updatedBy: string,
}