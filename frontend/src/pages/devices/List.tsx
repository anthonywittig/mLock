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
          }),
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
            <th scope="col">Warnings</th>
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
              <td>{renderEntityWarnings(entity)}</td>
              <td>
                <span className="btn">{renderEntityBatteryLevel(entity)}</span>
              </td>
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

  const renderEntityWarnings = (entity: DeviceT) => {
    const warnings = getOfflineWarnings(entity)
    warnings.push.apply(warnings, getLastRefreshedWarnings(entity))
    warnings.push.apply(warnings, getLastWentOfflineWarnings(entity))
    warnings.push.apply(warnings, getLockResponsivenessWarnings(entity))

    return <ListGroup>{warnings.map((warn) => warn)}</ListGroup>
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
      return <Badge bg="danger">Unknown</Badge>
    }

    if (level < 25) {
      return <Badge bg="danger">{level}%</Badge>
    }

    return <>{level}%</>
  }

  const getOfflineWarnings = (entity: DeviceT) => {
    const warnings: JSX.Element[] = []

    if (entity.rawDevice.status !== "ONLINE") {
      warnings.push(<ListGroup.Item variant="danger">Offline</ListGroup.Item>)
    }

    return warnings
  }

  const getLastRefreshedWarnings = (entity: DeviceT) => {
    const warnings: JSX.Element[] = []

    const lr = Date.parse(entity.lastRefreshedAt)
    const recently = sub(new Date(), { minutes: 70 })

    if (isBefore(lr, recently)) {
      const distance = formatDistance(lr, new Date(), { addSuffix: true })
      warnings.push(
        <ListGroup.Item variant="light">
          Last Data Sync: {distance}
        </ListGroup.Item>,
      )
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
      warnings.push(
        <ListGroup.Item variant="light">
          Went Offline: {distance}
        </ListGroup.Item>,
      )
    } else if (isAfter(lwond, recently)) {
      const distance = formatDistance(lwond, new Date(), { addSuffix: true })
      warnings.push(
        <ListGroup.Item variant="light">
          Went Online: {distance}
        </ListGroup.Item>,
      )
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
        warnings.push(
          <ListGroup.Item variant="light">
            The code {lc.code} was never added
          </ListGroup.Item>,
        )
        continue
      }
      if (lc.wasEnabledAt) {
        const wc = Date.parse(lc.wasEnabledAt)
        const minutesBetween = (wc - sa) / 1000 / 60
        if (expectedResponseInMinutes < minutesBetween) {
          const distance = formatDistance(sa, wc)
          warnings.push(
            <ListGroup.Item variant="light">
              Slow to Respond (took {distance} to add code {lc.code})
            </ListGroup.Item>,
          )
        } else {
          goodCode = true
        }
      } else {
        warnings.push(
          <ListGroup.Item variant="light">
            Not Responding (for code {lc.code})
          </ListGroup.Item>,
        )
      }
    }
  }

  return warnings
}

export { List, getLockResponsivenessWarnings }
