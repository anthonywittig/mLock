type AuditLogT = {
  id: string
  entries: AuditLogEntriesT[]
}

type AuditLogEntriesT = {
  createdAt: string
  log: string
}

type ClimateControlT = {
  climateControl: {
    id: string
    lastRefreshedAt: string
    desiredState: {
      endAt: string
      hvacMode: string
      note: string
      syncWithSettings: boolean
      temperature: number
    }
    rawClimateControl: {
      attributes: {
        current_temperature: number // actual
        friendly_name: string
        temperature: number // desired
      }
      entity_id: string
      state: string
    }
    syncWithReservations: boolean
  }
  unit: UnitT
}

type ClimateControlSetting = {
  hvacMode: string
  temperature: number
}

type DeviceLockCodeT = {
  code: string
  mode: string
  name: string
  slot: number
}

type DeviceT = {
  id: string
  unitId: string | null
  controllerId: string
  lastRefreshedAt: string
  lastWentOfflineAt: string | null
  lastWentOnlineAt: string | null
  managedLockCodes: DeviceManagedLockCodeT[]
  rawDevice: {
    battery: {
      batteryPowered: boolean
      level: number
    }
    categoryId: string
    lockCodes: [DeviceLockCode] | null
    name: string
    status: string
  }
}

type DeviceManagedLockCodeT = {
  id: string
  deviceId: string
  code: string
  note: string
  reservation: DeviceManagedLockCodeReservationT
  status: string
  startAt: string
  endAt: string
  startedAddingAt: string | null
  wasEnabledAt: string | null
  startedRemovingAt: string | null
  wasCompletedAt: string | null
}

type DeviceManagedLockCodeReservationT = {
  id: string
  sync: boolean
}

type UnitOccupancyStatusT = {
  date: string
  noon: {
    occupied: boolean
    managedLockCodes: DeviceManagedLockCodeT[]
  }
  fourPm: {
    occupied: boolean
    managedLockCodes: DeviceManagedLockCodeT[]
  }
}

type UnitT = {
  id: string
  name: string
  propertyId: string
  remotePropertyUrl: string
  updatedBy: string
}
