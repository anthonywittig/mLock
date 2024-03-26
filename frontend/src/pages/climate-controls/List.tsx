import React from "react"
import { Loading } from "../utils/Loading"
import { StandardFetch } from "../utils/FetchHelper"
import { formatDistance, isBefore, sub } from "date-fns"
import { Form, Button, ListGroup, Table } from "react-bootstrap"
import { Link } from "react-router-dom"

const Endpoint = "climate-controls"

const List = () => {
  const [entities, setEntities] = React.useState<ClimateControlT[]>([])
  const [climateControlOccupiedSettings, setClimateControlOccupiedSettings] =
    React.useState<ClimateControlSetting>({
      hvacMode: "off",
      temperature: 72,
    })
  const [climateControlVacantSettings, setClimateControlVacantSettings] =
    React.useState<ClimateControlSetting>({
      hvacMode: "off",
      temperature: 72,
    })
  const [loading, setLoading] = React.useState<boolean>(true)

  React.useEffect(() => {
    setLoading(true)

    StandardFetch(Endpoint, { method: "GET" })
      .then((response) => response.json())
      .then((response) => {
        setEntities(response.entities)
        setClimateControlOccupiedSettings(
          response.climateControlOccupiedSettings,
        )
        setClimateControlVacantSettings(response.climateControlVacantSettings)
        setLoading(false)
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }, [entities.length])

  const render = () => {
    return (
      <>
        <div className="card mb-2">
          <div className="card-header">
            <h2 className="card-title">Settings</h2>
          </div>
          <div className="card-body">{renderSettings()}</div>
        </div>
        <div className="card">
          <div className="card-header">
            <h2 className="card-title">Climate Controls</h2>
          </div>
          <div className="card-body">{renderEntities()}</div>
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
            <th scope="col">Unit</th>
            <th scope="col">Warnings</th>
            <th scope="col">Mode</th>
            <th scope="col">Set Temperature</th>
            <th scope="col">Actual Temperature</th>
            <th scope="col">Last Updated</th>
          </tr>
        </thead>
        <tbody>
          {entities.map((entity) => (
            <tr key={entity.climateControl.id}>
              <th scope="row">
                <Link to={"/climate-controls/" + entity.climateControl.id}>
                  <Button variant="link">
                    {entity.climateControl.rawClimateControl.attributes
                      .friendly_name ||
                      `${entity.climateControl.rawClimateControl.entity_id} (ID)`}
                  </Button>
                </Link>
              </th>
              <th>
                <Link to={"/units/" + entity.unit.id}>
                  <Button variant="link">{entity.unit.name}</Button>
                </Link>
              </th>
              <td>{renderStatus(entity)}</td>
              <td>{entity.climateControl.rawClimateControl.state}</td>
              <td>
                {entity.climateControl.rawClimateControl.attributes.temperature}
              </td>
              <td>
                {
                  entity.climateControl.rawClimateControl.attributes
                    .current_temperature
                }
              </td>
              <td>{renderLastRefreshedAt(entity)}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    )
  }

  const renderSettings = () => {
    if (loading) {
      return <Loading />
    }
    return (
      <>
        {[
          {
            name: "Occupied",
            settings: climateControlOccupiedSettings,
            settingsSetter: setClimateControlOccupiedSettings,
          },
          {
            name: "Vacant",
            settings: climateControlVacantSettings,
            settingsSetter: setClimateControlVacantSettings,
          },
        ].map((config) => (
          <>
            <h3>{config.name}</h3>
            <Form
              className="mb-3"
              onSubmit={(evt) => {
                evt.preventDefault()
                updateSettings()
              }}
            >
              <Form.Group
                controlId={`hvac-mode-${config.name}`}
                className="mb-3"
              >
                <Form.Label>HVAC Mode</Form.Label>
                <Form.Control
                  as="select"
                  onChange={(evt) =>
                    config.settingsSetter({
                      hvacMode: evt.target.value,
                      temperature: config.settings.temperature,
                    })
                  }
                >
                  <option
                    value="off"
                    selected={"off" === config.settings?.hvacMode}
                  >
                    Off
                  </option>
                  <option
                    value="cool"
                    selected={"cool" === config.settings?.hvacMode}
                  >
                    Cool
                  </option>
                  <option
                    value="heat"
                    selected={"heat" === config.settings?.hvacMode}
                  >
                    Heat
                  </option>
                </Form.Control>
              </Form.Group>

              <Form.Group
                controlId={`temperature-${config.name}`}
                className="mb-3"
              >
                <Form.Label>Temperature</Form.Label>
                <Form.Control
                  type="number"
                  value={config.settings?.temperature}
                  onChange={(evt) =>
                    config.settingsSetter({
                      hvacMode: config.settings.hvacMode,
                      temperature: Number(evt.target.value),
                    })
                  }
                />
              </Form.Group>

              <Button variant="secondary" type="submit">
                Update
              </Button>
            </Form>
          </>
        ))}
      </>
    )
  }

  const renderLastRefreshedAt = (entity: ClimateControlT) => {
    const lr = Date.parse(entity.climateControl.lastRefreshedAt)
    return formatDistance(lr, new Date(), { addSuffix: true })
  }

  const renderStatus = (entity: ClimateControlT) => {
    const warnings = getLastRefreshedWarnings(entity)
    return <ListGroup>{warnings.map((warn) => warn)}</ListGroup>
  }

  const getLastRefreshedWarnings = (entity: ClimateControlT) => {
    const warnings: JSX.Element[] = []

    const recently = sub(new Date(), { hours: 2 })
    const longAgo = sub(new Date(), { days: 1 })
    const lr = Date.parse(entity.climateControl.lastRefreshedAt)
    const distance = formatDistance(lr, new Date(), { addSuffix: true })

    if (isBefore(lr, longAgo)) {
      warnings.push(
        <ListGroup.Item variant="danger">
          Last Data Sync: {distance}
        </ListGroup.Item>,
      )
    } else if (isBefore(lr, recently)) {
      warnings.push(
        <ListGroup.Item variant="light">
          Last Data Sync: {distance}
        </ListGroup.Item>,
      )
    }

    return warnings
  }

  const updateSettings = () => {
    console.log("updateSettings")
    setLoading(true)

    StandardFetch(`${Endpoint}/settings`, {
      method: "PUT",
      body: JSON.stringify({
        climateControlOccupiedSettings,
        climateControlVacantSettings,
      }),
    })
      .then((response) => response.json())
      .then((response) => {
        setEntities(response.entities)
        setClimateControlOccupiedSettings(
          response.climateControlOccupiedSettings,
        )
        setClimateControlVacantSettings(response.climateControlVacantSettings)
        setLoading(false)
      })
      .catch((err) => {
        // TODO: indicate error.
        console.log(err)
      })
  }

  return render()
}

export { List }
