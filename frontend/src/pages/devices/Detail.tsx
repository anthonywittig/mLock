import React from "react"
import { Alert, Button, Form } from "react-bootstrap"
import { format, formatDistance, parseISO } from "date-fns"
import { useHistory, useRouteMatch } from "react-router-dom"
import { LockCode } from "./components/LockCode"
import { Loading } from "../utils/Loading"
import { StandardFetch } from "../utils/FetchHelper"

type MatchParams = { id: string }

const Endpoint = "devices"

export const Detail = () => {
  const [entity, setEntity] = React.useState<DeviceT>({
    id: "",
    unitId: "",
    lastRefreshedAt: "",
    lastWentOfflineAt: null,
    lastWentOnlineAt: null,
    managedLockCodes: [],
    rawDevice: {
      battery: {
        batteryPowered: false,
        level: 0,
      },
      categoryId: "",
      lockCodes: null,
      name: "",
      status: "",
    },
  })
  const [loading, setLoading] = React.useState<boolean>(true)
  const [auditLog, setAuditLog] = React.useState<AuditLogT>({
    id: "",
    entries: [],
  })
  const [units, setUnits] = React.useState<UnitT[]>([])
  const [unmanagedLockCodes, setUnmanagedLockCodes] = React.useState<
    DeviceLockCodeT[]
  >([])

  // `revision` is just to tell us when to pull the latest from the API.
  const [revision, setRevision] = React.useState<number>(0)
  const incrementRevision = () => {
    setRevision(revision + 1)
  }

  const m = useRouteMatch("/" + Endpoint + "/:id")
  const mp = m?.params as MatchParams
  const id = mp.id
  const history = useHistory()

  React.useEffect(() => {
    setLoading(true)

    StandardFetch(Endpoint + "/" + id, { method: "GET" })
      .then((response) => response.json())
      .then((response) => {
        setEntity(response.entity)
        setLoading(false)
        setAuditLog(response.extra.auditLog)
        setUnits(response.extra.units)
        setUnmanagedLockCodes(response.extra.unmanagedLockCodes)
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }, [id, revision])

  const detailFormUnitChange = (evt: React.ChangeEvent<HTMLSelectElement>) => {
    let val: string | null = evt.target.value
    if (val === "") {
      val = null
    }
    setEntity({
      ...entity,
      unitId: val,
    })
  }

  const formSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
    evt.preventDefault()

    setLoading(true)

    StandardFetch(Endpoint + "/" + id, {
      method: "PUT",
      body: JSON.stringify(entity),
    })
      .then((response) => response.json())
      .then((response) => {
        setEntity(response.entity)
        setLoading(false)
        history.push("/" + Endpoint + "/" + response.entity.id)
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }

  const render = () => {
    return (
      <>
        <div
          className="card"
          style={{ marginBottom: "1rem", marginTop: "1rem" }}
        >
          <div className="card-body">
            <h2 className="card-title">Details</h2>
            {renderEntity()}
          </div>
          {unmanagedLockCodes?.length ? (
            <div className="card-body">
              <h2 className="card-title">Unmanaged Lock Codes</h2>
              <Alert variant={"danger"}>
                These codes were added by another system.
              </Alert>
              <ul>
                {unmanagedLockCodes?.map((entry) => (
                  <li>
                    {entry.name} - {entry.code}
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <></>
          )}
          <div className="card-body">
            <h2 className="card-title">Managed Lock Codes</h2>
            {renderCurrentLockCodes()}
          </div>
          <div className="card-body">
            <h2 className="card-title">Add Managed Lock Code</h2>
            <LockCode
              deviceId={entity.id}
              managedLockCode={null}
              managedLockCodesUpdated={incrementRevision}
            />
          </div>
          <div className="card-body">
            <h2 className="card-title">Audit Log</h2>
            <div
              className="card"
              style={{ maxHeight: "300px", overflowY: "auto" }}
            >
              <ul className="list-group list-group-flush">
                {auditLog.entries.length ? (
                  auditLog.entries.map((entry) => (
                    <li className="list-group-item">
                      {format(parseISO(entry.createdAt), "L/d/yy h:mm aaa")} ---{" "}
                      {entry.log}
                    </li>
                  ))
                ) : (
                  <li>no entries yet</li>
                )}
              </ul>
            </div>
          </div>
        </div>
      </>
    )
  }

  const renderCurrentLockCodes = () => {
    if (loading) {
      return <Loading />
    }

    if (entity.managedLockCodes.length === 0) {
      return <p>There are no lock codes currently set.</p>
    }

    const getStatusValue = (status: string) => {
      switch (status) {
        case "Enabled":
          return 0
        case "Adding":
          return 1
        case "Removing":
          return 2
        case "Scheduled":
          return 3
        case "Complete":
          return 4
        default:
          console.log(`Couldn't identify status ${status}`)
          return -1
      }
    }

    entity.managedLockCodes.sort((a, b) => {
      const aValue = getStatusValue(a.status)
      const bValue = getStatusValue(b.status)

      const val = aValue - bValue
      if (val !== 0) {
        return val
      }

      if (aValue === 3) {
        return a.startAt.localeCompare(b.endAt)
      }

      if (aValue === 4) {
        return b.endAt.localeCompare(a.endAt)
      }

      return b.startAt.localeCompare(a.startAt)
    })

    return (
      <>
        {entity.managedLockCodes.map((lc) => {
          return (
            <div>
              <LockCode
                deviceId={entity.id}
                managedLockCode={lc}
                managedLockCodesUpdated={incrementRevision}
              />
              <br />
            </div>
          )
        })}
      </>
    )
  }

  const renderEntity = () => {
    if (loading) {
      return <Loading />
    }
    return (
      <Form onSubmit={(evt) => formSubmit(evt)}>
        <Form.Group>
          <Form.Label>Name</Form.Label>
          <Form.Control
            type="text"
            value={entity.rawDevice.name}
            disabled={true}
          />
        </Form.Group>

        <Form.Group>
          <Form.Label>Last Refreshed</Form.Label>
          <Form.Control
            type="text"
            value={formatDistance(
              Date.parse(entity.lastRefreshedAt),
              new Date(),
              { addSuffix: true }
            )}
            disabled={true}
          />
        </Form.Group>

        <Form.Group>
          <Form.Label>Status</Form.Label>
          <Form.Control
            type="text"
            value={entity.rawDevice.status}
            disabled={true}
          />
        </Form.Group>

        <Form.Group controlId="unit">
          <Form.Label>Unit</Form.Label>
          <Form.Control
            as="select"
            onChange={(evt) => detailFormUnitChange(evt as any)}
          >
            <option></option>
            {units.map((unit) => (
              <option value={unit.id} selected={unit.id === entity.unitId}>
                {unit.name}
              </option>
            ))}
          </Form.Control>
          <Form.Text className="text-muted">
            Lock codes will be created from the unit's reservations.
          </Form.Text>
        </Form.Group>

        <Button variant="secondary" type="submit">
          Update
        </Button>
      </Form>
    )
  }

  return render()
}
