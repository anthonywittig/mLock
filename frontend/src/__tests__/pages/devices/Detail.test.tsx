import { sortLockCodes } from "../../../pages/devices/Detail"

test("sort categories", () => {
  // We'll have two of each statuses, to help verify that they're actually sorted.
  const codes = sortLockCodes([
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Enabled",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Adding",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Removing",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Scheduled",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Complete",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Enabled",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Adding",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Removing",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Scheduled",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Complete",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
  ])

  expect(codes.length).toBe(10)
  expect(codes[0].status).toBe("Enabled")
  expect(codes[1].status).toBe("Enabled")
  expect(codes[2].status).toBe("Adding")
  expect(codes[3].status).toBe("Adding")
  expect(codes[4].status).toBe("Removing")
  expect(codes[5].status).toBe("Removing")
  expect(codes[6].status).toBe("Scheduled")
  expect(codes[7].status).toBe("Scheduled")
  expect(codes[8].status).toBe("Complete")
  expect(codes[9].status).toBe("Complete")
})

test("sort enabled", () => {
  const codes = sortLockCodes([
    {
      deviceId: "",
      code: "1",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Enabled",
      startAt: "2023-01-02T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "2",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Enabled",
      startAt: "2023-01-03T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "3",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Enabled",
      startAt: "2023-01-01T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
  ])

  expect(codes.length).toBe(3)
  expect(codes[0].code).toBe("2")
  expect(codes[1].code).toBe("1")
  expect(codes[2].code).toBe("3")
})

test("sort adding", () => {
  const codes = sortLockCodes([
    {
      deviceId: "",
      code: "1",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Adding",
      startAt: "2023-01-02T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "2",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Adding",
      startAt: "2023-01-03T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "3",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Adding",
      startAt: "2023-01-01T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
  ])

  expect(codes.length).toBe(3)
  expect(codes[0].code).toBe("2")
  expect(codes[1].code).toBe("1")
  expect(codes[2].code).toBe("3")
})

test("sort removing", () => {
  const codes = sortLockCodes([
    {
      deviceId: "",
      code: "1",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Removing",
      startAt: "2023-01-02T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "2",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Removing",
      startAt: "2023-01-03T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "3",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Removing",
      startAt: "2023-01-01T00:00:00Z",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
  ])

  expect(codes.length).toBe(3)
  expect(codes[0].code).toBe("2")
  expect(codes[1].code).toBe("1")
  expect(codes[2].code).toBe("3")
})

test("sort scheduled", () => {
  const nextYear = new Date().getFullYear() + 1
  const codes = sortLockCodes([
    {
      deviceId: "",
      code: "1",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Scheduled",
      startAt: `${nextYear}-01-02T00:00:00Z`,
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "2",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Scheduled",
      startAt: `${nextYear}-01-03T00:00:00Z`,
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "3",
      endAt: "",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Scheduled",
      startAt: `${nextYear}-01-01T00:00:00Z`,
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
  ])

  expect(codes.length).toBe(3)
  expect(codes[0].code).toBe("3")
  expect(codes[1].code).toBe("1")
  expect(codes[2].code).toBe("2")
})

test("sort complete", () => {
  const codes = sortLockCodes([
    {
      deviceId: "",
      code: "1",
      endAt: "2023-01-02T00:00:00Z",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Complete",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "2",
      endAt: "2023-01-03T00:00:00Z",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Complete",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
    {
      deviceId: "",
      code: "3",
      endAt: "2023-01-01T00:00:00Z",
      id: "",
      note: "",
      reservation: {
        id: "",
        sync: false,
      },
      status: "Complete",
      startAt: "",
      startedAddingAt: null,
      wasEnabledAt: null,
      startedRemovingAt: null,
      wasCompletedAt: null,
    },
  ])

  expect(codes.length).toBe(3)
  expect(codes[0].code).toBe("2")
  expect(codes[1].code).toBe("1")
  expect(codes[2].code).toBe("3")
})
