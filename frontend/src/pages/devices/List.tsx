import React from "react"
import {
  Badge,
  Button,
  ListGroup,
  OverlayTrigger,
  Table,
  Tooltip,
} from "react-bootstrap"
import { Loading } from "../utils/Loading"
import { StandardFetch } from "../utils/FetchHelper"
import { formatDistance, isAfter, isBefore, sub } from "date-fns"
import { Link } from "react-router-dom"

const Endpoint = "devices"

const List = () => {
  const [entities, setEntities] = React.useState<DeviceT[]>([])
  const [loading, setLoading] = React.useState<boolean>(true)
  const [units, setUnits] = React.useState<UnitT[]>([])

  React.useEffect(() => {
    setLoading(true)

    StandardFetch(Endpoint, { method: "GET" })
      .then((response) => response.json())
      .then((response) => {
        setEntities(response.entities)
        setLoading(false)
        setUnits(response.extra.units)
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }, [entities.length])

  const deleteDevice = (id: string) => {
    setLoading(true)

    StandardFetch(Endpoint + "/" + id, {
      method: "DELETE",
    })
      .then((_) => {
        setEntities(
          entities.filter((value) => {
            return value.id !== id
          })
        )
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }

  const render = () => {
    return (
      <>
        <div className="card">
          <div className="card-body">
            <h2 className="card-title">Devices</h2>
            {renderEntities()}
          </div>
        </div>
      </>
    )
  }

  const renderEntities = () => {
    if (loading) {
      return <Loading />
    }
    return (
      <Table responsive>
        <thead>
          <tr>
            <th scope="col">Name</th>
            <th scope="col">Status</th>
            <th scope="col">Battery</th>
            <th scope="col">Unit</th>
            <th scope="col">Actions</th>
          </tr>
        </thead>
        <tbody>
          {entities.map((entity) => (
            <tr key={entity.id}>
              <th scope="row">
                <Link to={"/devices/" + entity.id}>
                  <Button variant="link">{entity.rawDevice.name}</Button>
                </Link>
              </th>
              <td>{renderEntityStatus(entity)}</td>
              <td>{renderEntityBatteryLevel(entity)}</td>
              <td>
                <Link to={"/units/" + entity.unitId}>
                  <Button variant="link">
                    {units.find((e) => e.id === entity.unitId)?.name}
                  </Button>
                </Link>
              </td>
              <td>{renderDeleteButton(entity)}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    )
  }

  const renderDeleteButton = (entity: DeviceT) => {
    const lr = Date.parse(entity.lastRefreshedAt)
    // This should really be something much smaller, like 20 minutes, but since we have periods of time where we don't sync for an hour, we need something at least 60 minutes long.
    const recently = sub(new Date(), { minutes: 130 })

    if (isAfter(lr, recently)) {
      return (
        <OverlayTrigger
          overlay={
            <Tooltip id="tooltip-disabled">
              The device was recently pulled from the controller.
            </Tooltip>
          }
        >
          <span className="d-inline-block">
            <Button
              variant="secondary"
              disabled
              style={{ pointerEvents: "none" }}
            >
              Delete
            </Button>
          </span>
        </OverlayTrigger>
      )
    }

    return (
      <Button variant="secondary" onClick={() => deleteDevice(entity.id)}>
        Delete
      </Button>
    )
  }

  const renderEntityStatus = (entity: DeviceT) => {
    const warnings = getOfflineWarnings(entity)
    warnings.push.apply(warnings, getLastRefreshedWarnings(entity))
    warnings.push.apply(warnings, getLastWentOfflineWarnings(entity))
    warnings.push.apply(warnings, getLockResponsivenessWarnings(entity))

    return (
      <ListGroup className="flush">
        {warnings.map((warn) => (
          <ListGroup.Item className="border-0">{warn}</ListGroup.Item>
        ))}
      </ListGroup>
    )
  }

  const renderEntityBatteryLevel = (entity: DeviceT) => {
    if (!entity.rawDevice.battery.batteryPowered) {
      return <></>
    }

    const lu = entity.lastRefreshedAt
    const lud = Date.parse(lu)
    const recently = sub(new Date(), { days: 1, hours: 12 })
    const level = entity.rawDevice.battery.level

    if (isBefore(lud, recently) || level === null) {
      return <Badge variant="danger">Unknown</Badge>
    }

    if (level < 25) {
      return <Badge variant="danger">{level}%</Badge>
    }

    return <>{level}%</>
  }

  const getOfflineWarnings = (entity: DeviceT) => {
    const warnings: JSX.Element[] = []

    if (entity.rawDevice.status !== "ONLINE") {
      warnings.push(<Badge variant="danger">Offline</Badge>)
    }

    return warnings
  }

  const getLastRefreshedWarnings = (entity: DeviceT) => {
    const warnings: JSX.Element[] = []

    const lr = Date.parse(entity.lastRefreshedAt)
    const recently = sub(new Date(), { minutes: 70 })

    if (isBefore(lr, recently)) {
      const distance = formatDistance(lr, new Date(), { addSuffix: true })
      warnings.push(<>Last Data Sync: {distance}</>)
    }

    return warnings
  }

  const getLastWentOfflineWarnings = (entity: DeviceT) => {
    const warnings: JSX.Element[] = []

    const lwo = entity.lastWentOfflineAt
    if (lwo === null) {
      return warnings
    }

    const recently = sub(new Date(), { days: 1 })
    const lwond = Date.parse(entity.lastWentOnlineAt!)
    const lwoffd = Date.parse(entity.lastWentOfflineAt!)

    if (entity.rawDevice.status !== "ONLINE") {
      const distance = formatDistance(lwoffd, new Date(), { addSuffix: true })
      warnings.push(<>Went Offline: {distance}</>)
    } else if (isAfter(lwond, recently)) {
      const distance = formatDistance(lwond, new Date(), { addSuffix: true })
      warnings.push(<>Went Online: {distance}</>)
    }

    return warnings
  }

  return render()
}

const getLockResponsivenessWarnings = (entity: DeviceT) => {
  const warnings: JSX.Element[] = []

  if (entity.rawDevice.status !== "ONLINE") {
    return warnings
  }

  const tooSoon = sub(new Date(), { minutes: 10 })
  const expectedResponseInMinutes = 60
  let goodCode = false

  const sortedList = entity.managedLockCodes.sort((a, b) => {
    return Date.parse(b.startAt) - Date.parse(a.startAt)
  })

  for (let i = 0; i < sortedList.length; i++) {
    const lc = sortedList[i]

    if (
      (warnings.length && lc.status !== "Adding") ||
      (goodCode && lc.status !== "Adding") ||
      !lc.startedAddingAt //code is scheduled
    ) {
      continue
    }

    const sa = Date.parse(lc.startedAddingAt)
    if (isBefore(sa, tooSoon)) {
      if (lc.status === "Complete" && !lc.wasEnabledAt) {
        warnings.push(<>The code {lc.code} was never added</>)
        continue
      }
      if (lc.wasEnabledAt) {
        const wc = Date.parse(lc.wasEnabledAt)
        const minutesBetween = (wc - sa) / 1000 / 60
        if (expectedResponseInMinutes < minutesBetween) {
          const distance = formatDistance(sa, wc)
          warnings.push(
            <>
              Slow to Respond (took {distance} to add code {lc.code})
            </>
          )
        } else {
          goodCode = true
        }
      } else {
        warnings.push(<>Not Responding (for code {lc.code})</>)
      }
    }
  }

  return warnings
}

export { List, getLockResponsivenessWarnings }
