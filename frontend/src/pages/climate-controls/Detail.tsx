import React from "react"
import { Button, Form } from "react-bootstrap"
import { format, formatDistance, parseISO } from "date-fns"
import { useNavigate, useMatch } from "react-router-dom"
import { Loading } from "../utils/Loading"
import { StandardFetch } from "../utils/FetchHelper"

type MatchParams = { id: string }

const Endpoint = "climate-controls"

const Detail = () => {
  const [auditLog, setAuditLog] = React.useState<AuditLogT>({
    id: "",
    entries: [],
  })
  const [entity, setEntity] = React.useState<ClimateControlT>({
    climateControl: {
      id: "",
      lastRefreshedAt: "",
      desiredState: {
        abandonAfter: "",
        endAt: "",
        hvacMode: "",
        note: "",
        syncWithSettings: false,
        temperature: -1,
        wasSuccessfulAt: null,
      },
      rawClimateControl: {
        attributes: {
          current_temperature: -1,
          friendly_name: "",
          temperature: -1,
        },
        entity_id: "",
        state: "",
      },
      syncWithReservations: false,
    },
    unit: {
      id: "",
      name: "",
      propertyId: "",
      remotePropertyUrl: "",
      updatedBy: "",
    },
  })
  const [loading, setLoading] = React.useState<boolean>(true)
  const [syncWithReservations, setSyncWithReservations] =
    React.useState<boolean>(false)
  const [unitOccupancyStatuses, setUnitOccupancyStatuses] = React.useState<
    UnitOccupancyStatusT[]
  >([])

  const m = useMatch(Endpoint + "/:id")
  const mp = m?.params as MatchParams
  const id = mp.id
  const navigate = useNavigate()

  React.useEffect(() => {
    setLoading(true)

    StandardFetch(Endpoint + "/" + id, { method: "GET" })
      .then((response) => response.json())
      .then((response) => {
        setAuditLog(response.extra.auditLog)
        setEntity(response.entity)
        setSyncWithReservations(
          response.entity.climateControl.syncWithReservations,
        )
        setUnitOccupancyStatuses(response.extra.unitOccupancyStatuses)
        setLoading(false)
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }, [id])

  const formSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
    evt.preventDefault()

    setLoading(true)

    StandardFetch(Endpoint + "/" + id, {
      method: "PUT",
      body: JSON.stringify({
        syncWithReservations: syncWithReservations,
      }),
    })
      .then((response) => response.json())
      .then((response) => {
        setAuditLog(response.extra.auditLog)
        setEntity(response.entity)
        setSyncWithReservations(
          response.entity.climateControl.syncWithReservations,
        )
        setUnitOccupancyStatuses(response.extra.unitOccupancyStatuses)
        setLoading(false)
        navigate("/" + Endpoint + "/" + response.entity.climateControl.id)
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
          <div className="card-body">
            <h2 className="card-title">Active Reservations At</h2>
            {renderOccupancyStatuses()}
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
                  <li>&nbsp;no entries yet</li>
                )}
              </ul>
            </div>
          </div>
        </div>
      </>
    )
  }

  const renderEntity = () => {
    if (loading) {
      return <Loading />
    }
    return (
      <Form onSubmit={(evt) => formSubmit(evt)}>
        <Form.Group className="mb-3">
          <Form.Label>Name</Form.Label>
          <Form.Control
            type="text"
            value={
              entity.climateControl.rawClimateControl.attributes.friendly_name
            }
            disabled={true}
          />
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>Last Refreshed</Form.Label>
          <Form.Control
            type="text"
            value={formatDistance(
              Date.parse(entity.climateControl.lastRefreshedAt),
              new Date(),
              { addSuffix: true },
            )}
            disabled={true}
          />
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>HVAC Mode</Form.Label>
          <Form.Control
            type="text"
            value={entity.climateControl.rawClimateControl.state}
            disabled={true}
          />
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>Temperature</Form.Label>
          <Form.Control
            type="text"
            value={
              entity.climateControl.rawClimateControl.attributes.temperature
            }
            disabled={true}
          />
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>Desired Was Set At</Form.Label>
          {entity.climateControl.desiredState.wasSuccessfulAt ? (
            <Form.Control
              type="datetime-local"
              defaultValue={format(
                parseISO(entity.climateControl.desiredState.wasSuccessfulAt),
                "yyyy-MM-dd'T'HH:mm",
              )}
              disabled={true}
            />
          ) : (
            <Form.Control type="text" defaultValue="Not Yet" disabled={true} />
          )}
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>Abandon Desired At</Form.Label>
          <Form.Control
            type="datetime-local"
            defaultValue={format(
              parseISO(entity.climateControl.desiredState.abandonAfter),
              "yyyy-MM-dd'T'HH:mm",
            )}
            disabled={true}
          />
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>Desired HVAC Mode</Form.Label>
          <Form.Control
            type="text"
            value={entity.climateControl.desiredState.hvacMode}
            disabled={true}
          />
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>Desired Temperature</Form.Label>
          <Form.Control
            type="text"
            value={entity.climateControl.desiredState.temperature}
            disabled={true}
          />
        </Form.Group>

        <Form.Group className="mb-3" controlId="formBasicCheckbox">
          <Form.Check
            type="checkbox"
            label="Sync With Reservation"
            checked={syncWithReservations}
            onChange={(e) => setSyncWithReservations(e.target.checked)}
          />
          <small id="formBasicCheckboxHelp" className="form-text text-muted">
            Technically, this syncs with lock codes that are associated with a
            reservation.
          </small>
        </Form.Group>

        <Button variant="secondary" type="submit">
          Update
        </Button>
      </Form>
    )
  }

  const renderOccupancyStatuses = () => {
    if (loading) {
      return <Loading />
    }
    return (
      <table className="table table-responsive-sm">
        <thead>
          <tr>
            <th scope="col">Date</th>
            <th scope="col">Noon</th>
            <th scope="col">4 PM</th>
          </tr>
        </thead>
        <tbody>
          {unitOccupancyStatuses.map((status) => (
            <tr>
              <th scope="row">{format(parseISO(status.date), "LL/dd/yyyy")}</th>
              <td>
                {status.noon.occupied
                  ? status.noon.managedLockCodes[0].reservation.id
                  : "-"}
              </td>
              <td>
                {status.fourPm.occupied
                  ? status.fourPm.managedLockCodes[0].reservation.id
                  : "-"}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    )
  }

  return render()
}

export { Detail }
