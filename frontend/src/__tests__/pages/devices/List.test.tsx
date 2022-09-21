import { getLockResponsivenessWarnings } from '../../../pages/devices/List'

test('warning for single none-responsive code', () => {
    const warnings = getLockResponsivenessWarnings({
        id: "",
        unitId: "",
        lastRefreshedAt: "",
        lastWentOfflineAt: null,
        lastWentOnlineAt: null,
        managedLockCodes: [
            {
                deviceId: "",
                code: "1234",
                endAt: "2030-01-01T01:00:00Z",
                id: "a25cd2cd-bf1b-4193-9c71-1d796d47aae7",
                note: "Lock code present.",
                reservation: {
                    id: "",
                    sync: false
                },
                status: "Enabled",
                startAt: "2020-01-01T01:00:00Z",
                startedAddingAt: "2020-01-01T01:00:00Z",
                wasEnabledAt: "2020-01-02T01:00:00Z",
                startedRemovingAt: null,
                wasCompletedAt: null
            },
        ],
        rawDevice: {
            battery: {
                batteryPowered: true,
                level: 0,
            },
            categoryId: "",
            lockCodes: null,
            name: "",
            status: "ONLINE",
        }
    })

    expect(warnings.length).toBe(1)
    expect(warnings[0]).toStrictEqual(<>Slow to Respond (took {"1 day"} to add code {"1234"})</>)
})
