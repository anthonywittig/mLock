import { getLockResponsivenessWarnings } from "../../../pages/devices/List"

test("scheduled is ignored", () => {
  const warnings = getLockResponsivenessWarnings({
    id: "",
    unitId: "",
    lastRefreshedAt: "",
    lastWentOfflineAt: null,
    lastWentOnlineAt: null,
    managedLockCodes: [
      {
        deviceId: "",
        code: "0000",
        endAt: "2022-10-13T01:00:00Z",
        id: "edfb9b4a-7ce8-4b6f-8dc3-8f0f30d511a8",
        note: "Added by awesomeiscooltoo@gmail.com.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Scheduled",
        startAt: "2022-10-11T01:44:00Z",
        startedAddingAt: null,
        wasEnabledAt: null,
        startedRemovingAt: null,
        wasCompletedAt: null,
      },
      {
        deviceId: "",
        code: "9999",
        endAt: "2021-09-15T11:30:00-06:00",
        id: "81e8be93-6795-41ea-94d9-df76df002750",
        note: "Code was removed.",
        reservation: {
          id: "8833936@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2021-09-13T15:00:00-06:00",
        startedAddingAt: "2021-09-13T00:00:00Z",
        wasEnabledAt: "2021-09-14T00:00:00Z",
        startedRemovingAt: "2021-09-15T00:00:00Z",
        wasCompletedAt: "2021-09-15T00:00:00Z",
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
    },
  })

  expect(warnings.length).toBe(1)
  expect(warnings[0]).toStrictEqual(
    <>
      Slow to Respond (took {"1 day"} to add code {"9999"})
    </>
  )
})

test("Warning for older adding code", () => {
  const warnings = getLockResponsivenessWarnings({
    id: "",
    unitId: "",
    lastRefreshedAt: "",
    lastWentOfflineAt: null,
    lastWentOnlineAt: null,
    managedLockCodes: [
      {
        deviceId: "",
        code: "0000",
        endAt: "2021-09-15T11:30:00-06:00",
        id: "81e8be93-6795-41ea-94d9-df76df002750",
        note: "Code was removed.",
        reservation: {
          id: "8833936@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2021-09-13T15:00:00-06:00",
        startedAddingAt: "2021-09-13T00:00:00Z",
        wasEnabledAt: "2021-09-13T00:00:00Z",
        startedRemovingAt: "2021-09-15T00:00:00Z",
        wasCompletedAt: "2021-09-15T00:00:00Z",
      },
      {
        deviceId: "",
        code: "9999",
        endAt: "2022-09-15T11:30:00-06:00",
        id: "91e8be93-6795-41ea-94d9-df76df002750",
        note: "",
        reservation: {
          id: "8833936@LiveRez.com",
          sync: true,
        },
        status: "Adding",
        startAt: "2020-09-13T15:00:00-06:00",
        startedAddingAt: "2020-09-13T00:00:00Z",
        wasEnabledAt: null,
        startedRemovingAt: null,
        wasCompletedAt: null,
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
    },
  })

  expect(warnings.length).toBe(1)
  expect(warnings[0]).toStrictEqual(<>Not Responding (for code {"9999"})</>)
})

test("warning for single none-responsive code", () => {
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
          sync: false,
        },
        status: "Enabled",
        startAt: "2020-01-01T01:00:00Z",
        startedAddingAt: "2020-01-01T01:00:00Z",
        wasEnabledAt: "2020-01-02T01:00:00Z",
        startedRemovingAt: null,
        wasCompletedAt: null,
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
    },
  })

  expect(warnings.length).toBe(1)
  expect(warnings[0]).toStrictEqual(
    <>
      Slow to Respond (took {"1 day"} to add code {"1234"})
    </>
  )
})

test("warning for single none-responsive completed code", () => {
  const warnings = getLockResponsivenessWarnings({
    id: "",
    unitId: "",
    lastRefreshedAt: "",
    lastWentOfflineAt: null,
    lastWentOnlineAt: null,
    managedLockCodes: [
      {
        deviceId: "",
        code: "3936",
        endAt: "2022-09-15T11:30:00-06:00",
        id: "81e8be93-6795-41ea-94d9-df76df002750",
        note: "Code was removed.",
        reservation: {
          id: "8833936@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-09-13T15:00:00-06:00",
        startedAddingAt: "2022-09-13T21:01:35.24791605Z",
        wasEnabledAt: "2022-09-14T21:01:35.24791605Z",
        startedRemovingAt: "2022-09-15T18:11:12.388932124Z",
        wasCompletedAt: "2022-09-15T19:11:20.461120924Z",
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
    },
  })

  expect(warnings.length).toBe(1)
  expect(warnings[0]).toStrictEqual(
    <>
      Slow to Respond (took {"1 day"} to add code {"3936"})
    </>
  )
})

test("warning for single never added completed code", () => {
  const warnings = getLockResponsivenessWarnings({
    id: "",
    unitId: "",
    lastRefreshedAt: "",
    lastWentOfflineAt: null,
    lastWentOnlineAt: null,
    managedLockCodes: [
      {
        deviceId: "",
        code: "4622",
        endAt: "2022-10-12T02:10:00Z",
        id: "224818c5-735c-4911-b0ab-24c43c17d45e",
        note: "Code was removed.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Complete",
        startAt: "2022-06-15T22:15:25.408Z",
        startedAddingAt: "2022-10-12T01:36:00.600297352Z",
        wasEnabledAt: null,
        startedRemovingAt: null,
        wasCompletedAt: "2022-10-12T02:10:59.840451195Z",
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
    },
  })

  expect(warnings.length).toBe(1)
  expect(warnings[0]).toStrictEqual(<>The code {"4622"} was never added</>)
})

test("single warning for adding code (ignore an older completed code that never added)", () => {
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
        note: "Attempting to add lock code.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Adding",
        startAt: "2021-01-01T01:00:00Z",
        startedAddingAt: "2021-01-01T01:00:00Z",
        wasEnabledAt: null,
        startedRemovingAt: null,
        wasCompletedAt: null,
      },
      {
        deviceId: "",
        code: "4622",
        endAt: "2020-01-02T01:00:00Z",
        id: "224818c5-735c-4911-b0ab-24c43c17d45e",
        note: "Code was removed.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Complete",
        startAt: "2020-01-01T01:00:00Z",
        startedAddingAt: "2020-01-01T01:00:00Z",
        wasEnabledAt: null,
        startedRemovingAt: null,
        wasCompletedAt: "2020-01-02T01:00:00Z",
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
    },
  })

  expect(warnings.length).toBe(1)
  expect(warnings[0]).toStrictEqual(<>Not Responding (for code {"1234"})</>)
})

test("no warning for a quickly added code (ignore an older completed code that never added)", () => {
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
          sync: false,
        },
        status: "Enabled",
        startAt: "2021-01-01T01:00:00Z",
        startedAddingAt: "2021-01-01T01:00:00Z",
        wasEnabledAt: "2021-01-01T01:00:00Z",
        startedRemovingAt: null,
        wasCompletedAt: null,
      },
      {
        deviceId: "",
        code: "4622",
        endAt: "2020-01-02T01:00:00Z",
        id: "224818c5-735c-4911-b0ab-24c43c17d45e",
        note: "Code was removed.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Complete",
        startAt: "2020-01-01T01:00:00Z",
        startedAddingAt: "2020-01-01T01:00:00Z",
        wasEnabledAt: null,
        startedRemovingAt: null,
        wasCompletedAt: "2020-01-02T01:00:00Z",
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
    },
  })

  expect(warnings.length).toBe(0)
})

test("warning for a completed code despite having an old enabled code", () => {
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
          sync: false,
        },
        status: "Enabled",
        startAt: "2020-01-01T01:00:00Z",
        startedAddingAt: "2020-01-01T01:00:00Z",
        wasEnabledAt: "2020-01-01T01:00:00Z",
        startedRemovingAt: null,
        wasCompletedAt: null,
      },
      {
        deviceId: "",
        code: "4622",
        endAt: "2022-01-02T01:00:00Z",
        id: "224818c5-735c-4911-b0ab-24c43c17d45e",
        note: "Code was removed.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Complete",
        startAt: "2022-01-01T01:00:00Z",
        startedAddingAt: "2022-01-01T01:00:00Z",
        wasEnabledAt: null,
        startedRemovingAt: null,
        wasCompletedAt: "2022-01-02T01:00:00Z",
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
    },
  })

  expect(warnings.length).toBe(1)
  expect(warnings[0]).toStrictEqual(<>The code {"4622"} was never added</>)
})

test("no warning for single responsive code", () => {
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
          sync: false,
        },
        status: "Enabled",
        startAt: "2020-01-01T01:00:00Z",
        startedAddingAt: "2020-01-01T01:00:00Z",
        wasEnabledAt: "2020-01-01T01:00:00Z",
        startedRemovingAt: null,
        wasCompletedAt: null,
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
    },
  })

  expect(warnings.length).toBe(0)
})

test("jacob will give me a name", () => {
  const warnings = getLockResponsivenessWarnings({
    id: "",
    unitId: "",
    lastRefreshedAt: "",
    lastWentOfflineAt: null,
    lastWentOnlineAt: null,
    managedLockCodes: [
      {
        deviceId: "",
        code: "4622",
        endAt: "2026-06-23T21:00:00Z",
        id: "a25cd2cd-bf1b-4193-9c71-1d796d47aae7",
        note: "Lock code present.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Enabled",
        startAt: "2022-06-22T21:57:10.987Z",
        startedAddingAt: "2022-06-22T22:01:07.075767618Z",
        wasEnabledAt: "2022-06-22T22:09:06.716284723Z",
        startedRemovingAt: null,
        wasCompletedAt: null,
      },
      {
        deviceId: "",
        code: "0205",
        endAt: "2026-06-23T21:00:00Z",
        id: "739bb535-d87a-4e76-9fba-aa7da75f2716",
        note: "Lock code present.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Enabled",
        startAt: "2022-06-22T21:57:10.987Z",
        startedAddingAt: "2022-06-22T22:01:07.225823329Z",
        wasEnabledAt: "2022-06-22T22:04:44.007053851Z",
        startedRemovingAt: null,
        wasCompletedAt: null,
      },
      {
        deviceId: "",
        code: "7598",
        endAt: "2022-09-19T11:30:00-06:00",
        id: "858dc6bb-9fb5-4741-bb95-1d6e9471028c",
        note: "Lock code present.",
        reservation: {
          id: "8837598@LiveRez.com",
          sync: true,
        },
        status: "Enabled",
        startAt: "2022-09-15T15:00:00-06:00",
        startedAddingAt: "2022-09-15T21:01:24.029599197Z",
        wasEnabledAt: "2022-09-15T21:11:30.249794228Z",
        startedRemovingAt: null,
        wasCompletedAt: null,
      },
      {
        deviceId: "",
        code: "3936",
        endAt: "2022-09-15T11:30:00-06:00",
        id: "81e8be93-6795-41ea-94d9-df76df002750",
        note: "Code was removed.",
        reservation: {
          id: "8833936@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-09-13T15:00:00-06:00",
        startedAddingAt: "2022-09-13T21:01:35.24791605Z",
        wasEnabledAt: "2022-09-13T21:06:38.336081606Z",
        startedRemovingAt: "2022-09-15T18:11:12.388932124Z",
        wasCompletedAt: "2022-09-15T19:11:20.461120924Z",
      },
      {
        deviceId: "",
        code: "8226",
        endAt: "2022-09-11T11:30:00-06:00",
        id: "a44c4514-e18d-4b76-8fd5-8b06ac65ab4b",
        note: "Code was removed.",
        reservation: {
          id: "8848226@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-09-10T15:00:00-06:00",
        startedAddingAt: "2022-09-10T21:01:04.153359782Z",
        wasEnabledAt: "2022-09-11T03:51:02.630865027Z",
        startedRemovingAt: "2022-09-11T18:10:58.098390464Z",
        wasCompletedAt: "2022-09-11T19:11:04.678902041Z",
      },
      {
        deviceId: "",
        code: "5509",
        endAt: "2022-09-10T11:30:00-06:00",
        id: "021b02d7-d039-4a52-9ac4-66520290260b",
        note: "Code was removed.",
        reservation: {
          id: "8715509@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-09-09T15:00:00-06:00",
        startedAddingAt: "2022-09-09T21:01:04.596072011Z",
        wasEnabledAt: "2022-09-09T21:06:01.43510538Z",
        startedRemovingAt: "2022-09-10T18:10:59.372011808Z",
        wasCompletedAt: "2022-09-10T19:11:05.032618864Z",
      },
      {
        deviceId: "",
        code: "7906",
        endAt: "2022-09-05T11:30:00-06:00",
        id: "7e6433f5-e0dc-4867-a81c-2b7fb02e6a89",
        note: "Code was removed.",
        reservation: {
          id: "8837906@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-09-04T15:00:00-06:00",
        startedAddingAt: "2022-09-04T21:01:01.327958779Z",
        wasEnabledAt: "2022-09-04T21:06:03.105457589Z",
        startedRemovingAt: "2022-09-05T18:10:57.692747479Z",
        wasCompletedAt: "2022-09-05T19:10:57.474523414Z",
      },
      {
        deviceId: "",
        code: "0603",
        endAt: "2022-09-04T11:30:00-06:00",
        id: "dd051d77-f64e-4988-91f4-e934875ae263",
        note: "Code was removed.",
        reservation: {
          id: "8810603@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-09-03T15:00:00-06:00",
        startedAddingAt: "2022-09-03T21:01:00.470735128Z",
        wasEnabledAt: "2022-09-03T21:05:59.461165465Z",
        startedRemovingAt: "2022-09-04T18:11:01.246584434Z",
        wasCompletedAt: "2022-09-04T19:11:02.055162664Z",
      },
      {
        deviceId: "",
        code: "1916",
        endAt: "2022-08-27T11:30:00-06:00",
        id: "7ea2bd3c-5f6f-4cf2-8293-28d80bf500d3",
        note: "Code was removed.",
        reservation: {
          id: "8821916@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-08-26T15:30:00-06:00",
        startedAddingAt: "2022-08-26T21:30:31.770493777Z",
        wasEnabledAt: "2022-08-26T21:35:34.808132057Z",
        startedRemovingAt: "2022-08-27T17:30:34.549278941Z",
        wasCompletedAt: "2022-08-27T17:35:32.151276248Z",
      },
      {
        deviceId: "",
        code: "1111",
        endAt: "2022-08-25T23:30:00Z",
        id: "92b253b0-328c-4350-a80a-97b7ae8314d6",
        note: "Code was removed.",
        reservation: {
          id: "",
          sync: false,
        },
        status: "Complete",
        startAt: "2022-08-25T22:37:24.017Z",
        startedAddingAt: "2022-08-25T22:40:04.93346993Z",
        wasEnabledAt: "2022-08-25T22:44:01.437214218Z",
        startedRemovingAt: "2022-08-25T23:32:03.854830585Z",
        wasCompletedAt: "2022-08-25T23:36:01.264507644Z",
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
    },
  })
  //we want this to have no warnings
  expect(warnings.length).toBe(0)
})

test("no warning for old code", () => {
  const warnings = getLockResponsivenessWarnings({
    id: "",
    unitId: "",
    lastRefreshedAt: "",
    lastWentOfflineAt: null,
    lastWentOnlineAt: null,
    managedLockCodes: [
      {
        deviceId: "",
        code: "0000",
        endAt: "2022-09-15T11:30:00-06:00",
        id: "81e8be93-6795-41ea-94d9-df76df002750",
        note: "Code was removed.",
        reservation: {
          id: "8833936@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2022-09-13T00:00:00-06:00",
        startedAddingAt: "2022-09-13T00:00:00Z",
        wasEnabledAt: "2022-09-13T00:00:00Z",
        startedRemovingAt: "2022-09-15T00:00:00Z",
        wasCompletedAt: "2022-09-15T00:00:00Z",
      },
      {
        deviceId: "",
        code: "9999",
        endAt: "2021-09-15T11:30:00-06:00",
        id: "81e8be93-6795-41ea-94d9-df76df002750",
        note: "Code was removed.",
        reservation: {
          id: "8833936@LiveRez.com",
          sync: true,
        },
        status: "Complete",
        startAt: "2021-09-13T15:00:00-06:00",
        startedAddingAt: "2021-09-13T00:00:00Z",
        wasEnabledAt: "2021-09-14T00:00:00Z",
        startedRemovingAt: "2021-09-15T00:00:00Z",
        wasCompletedAt: "2021-09-15T00:00:00Z",
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
    },
  })
  expect(warnings.length).toBe(0)
})
