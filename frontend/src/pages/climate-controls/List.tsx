import React from "react"
import { Table } from "react-bootstrap"
import { Loading } from "../utils/Loading"
import { StandardFetch } from "../utils/FetchHelper"
import { formatDistance } from "date-fns"

const Endpoint = "climate-controls"

const List = () => {
  const [entities, setEntities] = React.useState<ClimateControlT[]>([])
  const [loading, setLoading] = React.useState<boolean>(true)

  React.useEffect(() => {
    setLoading(true)

    StandardFetch(Endpoint, { method: "GET" })
      .then((response) => response.json())
      .then((response) => {
        setEntities(response.entities)
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
        <div className="card">
          <div className="card-body">
            <h2 className="card-title">Climate Controls</h2>
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
            <th scope="col">Set Temperature</th>
            <th scope="col">Actual Temperature</th>
            <th scope="col">Last Updated</th>
          </tr>
        </thead>
        <tbody>
          {entities.map((entity) => (
            <tr key={entity.id}>
              <th scope="row">
                {entity.rawClimateControl.attributes.friendly_name}
              </th>
              <td>{entity.rawClimateControl.state}</td>
              <td>{entity.rawClimateControl.attributes.temperature}</td>
              <td>{entity.rawClimateControl.attributes.current_temperature}</td>
              <td>{renderLastRefreshedAt(entity)}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    )
  }

  const renderLastRefreshedAt = (entity: ClimateControlT) => {
    const lr = Date.parse(entity.lastRefreshedAt)
    return formatDistance(lr, new Date(), { addSuffix: true })
  }

  return render()
}

export { List }
