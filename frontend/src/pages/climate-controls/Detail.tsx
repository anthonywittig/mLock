import React from "react"
import { Button, Form } from "react-bootstrap"
import { format, formatDistance, parseISO } from "date-fns"
import { useNavigate, useMatch } from "react-router-dom"
import { Loading } from "../utils/Loading"
import { StandardFetch } from "../utils/FetchHelper"

type MatchParams = { id: string }

const Endpoint = "climate-controls"

const Detail = () => {
  const [entity, setEntity] = React.useState<ClimateControlT>({
    climateControl: {
      id: "",
      lastRefreshedAt: "",
      rawClimateControl: {
        attributes: {
          current_temperature: -1,
          friendly_name: "",
          temperature: -1,
        },
        entity_id: "",
        state: "",
      },
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
  const [auditLog, setAuditLog] = React.useState<AuditLogT>({
    id: "",
    entries: [],
  })

  const m = useMatch(Endpoint + "/:id")
  const mp = m?.params as MatchParams
  const id = mp.id
  const navigate = useNavigate()

  React.useEffect(() => {
    setLoading(true)

    StandardFetch(Endpoint + "/" + id, { method: "GET" })
      .then((response) => response.json())
      .then((response) => {
        setEntity(response.entity)
        setLoading(false)
        setAuditLog(response.extra.auditLog)
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
      body: JSON.stringify(entity),
    })
      .then((response) => response.json())
      .then((response) => {
        setEntity(response.entity)
        setLoading(false)
        navigate("/" + Endpoint + "/" + response.entity.id)
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

        <Button variant="secondary" type="submit">
          Update
        </Button>
      </Form>
    )
  }

  return render()
}

export { Detail }
