import React from "react"
import { Button } from "react-bootstrap"
import { Link, useHistory } from "react-router-dom"
import { StandardFetch } from "../utils/FetchHelper"
import { History } from "history"

type Adder = (
  id: string,
  name: string,
  propertyId: string,
  updatedBy: string,
) => void
type IdAction = (id: string) => void

type Device = {
  id: string
  name: string
}

type Property = {
  id: string
  name: string
  updatedBy: string
}

type Props = {
  devices: Device[]
  entityId: string
  entityName: string
  propertyId: string
  updatedBy: string
  properties: Property[]
  addEntity: Adder | null
  removeEntity: IdAction | null
}

type State = {
  entityFieldsDisabled: boolean
  setEntityFieldsDisabled: React.Dispatch<React.SetStateAction<boolean>>
  entityName: string
  setEntityName: React.Dispatch<React.SetStateAction<string>>
  propertyId: string
  setPropertyId: React.Dispatch<React.SetStateAction<string>>
  entityState: string
  setEntityState: React.Dispatch<React.SetStateAction<string>>
  history: History
}

const Endpoint = "units"

export const Unit = (props: Props) => {
  const state = GetState(props)
  return render(props, state)
}

function GetState(props: Props): State {
  const [entityFieldsDisabled, setEntityFieldsDisabled] =
    React.useState<boolean>(false)
  const [entityName, setEntityName] = React.useState<string>(props.entityName)
  const [propertyId, setPropertyId] = React.useState<string>(props.propertyId)
  const [entityState, setEntityState] = React.useState<string>(
    props.entityName ? "exists" : "new",
  )
  const history = useHistory()
  return {
    entityFieldsDisabled,
    setEntityFieldsDisabled,
    entityName,
    setEntityName,
    propertyId,
    setPropertyId,
    entityState,
    setEntityState,
    history,
  }
}

function resetState(props: Props, state: State) {
  // Some dupliation from `getState`...
  state.setEntityFieldsDisabled(false)
  state.setEntityName(props.entityName)
  state.setPropertyId(props.propertyId)
  state.setEntityState(props.entityId ? "exists" : "new")
}

function removeClick(props: Props, id: string) {
  StandardFetch(Endpoint + "/" + id, { method: "DELETE" })
    .then((response) => {
      if (response.status === 200) {
        if (props.removeEntity) {
          props.removeEntity(id)
        } else {
          throw new Error("removeEntry is null")
        }
      }
    })
    .catch((err) => {
      // TODO: need to indicate error.
      console.log("error: " + err)
    })
}

function nameClick(state: State, id: string) {
  state.history.push("/units/" + id)
}

function updateEntityName(
  state: State,
  evt: React.ChangeEvent<HTMLInputElement>,
) {
  state.setEntityName(evt.target.value)
}

function updatePropertyId(
  state: State,
  evt: React.ChangeEvent<HTMLSelectElement>,
) {
  state.setPropertyId(evt.target.value)
}

function newEntitySubmit(props: Props, state: State) {
  if (!state.entityName || !state.propertyId) {
    // TODO: indicate error.
    return
  }

  state.setEntityFieldsDisabled(true)

  StandardFetch(Endpoint, {
    method: "POST",
    body: JSON.stringify({
      name: state.entityName,
      propertyId: state.propertyId,
    }),
  })
    .then((response) => response.json())
    .then((response) => {
      // add to parent
      let e = response.entity
      if (props.addEntity) {
        props.addEntity(e.id, e.name, e.propertyId, e.updatedBy)
        resetState(props, state)
      } else {
        throw new Error("addEntity is null")
      }
    })
    .catch((err) => {
      // TODO: indicate error.
      state.setEntityFieldsDisabled(false)
    })
}

function render(props: Props, state: State) {
  if (state.entityState === "new") {
    return (
      <tr key="newEntity">
        <th scope="row">
          <input
            type="text"
            className="form-control"
            id="newName"
            placeholder="Name"
            value={state.entityName}
            onChange={(evt) => updateEntityName(state, evt)}
            disabled={state.entityFieldsDisabled}
            onKeyUp={(evt) =>
              evt.key === "Enter" ? newEntitySubmit(props, state) : ""
            }
          />
        </th>
        <td></td>
        <td>
          <select
            id="newProperty"
            className="form-control"
            onChange={(evt) => updatePropertyId(state, evt)}
            disabled={state.entityFieldsDisabled}
          >
            <option></option>
            {props.properties.map((property) => (
              <option
                value={property.id}
                selected={property.id === state.propertyId}
              >
                {property.name}
              </option>
            ))}
          </select>
        </td>
        <td>
          <Button
            variant="secondary"
            onClick={() => newEntitySubmit(props, state)}
            disabled={state.entityFieldsDisabled}
          >
            Create
          </Button>
        </td>
      </tr>
    )
  }

  let devices = <></>
  if (props.devices.length === 1) {
    devices = (
      <Link to={"/devices/" + props.devices[0].id}>
        <Button variant="link">{props.devices[0].name}</Button>
      </Link>
    )
  } else {
    devices = (
      <ul>
        {props.devices.map((device) => (
          <li>
            <Link to={"/devices/" + device.id}>
              <Button variant="link">{device.name}</Button>
            </Link>
          </li>
        ))}
      </ul>
    )
  }

  return (
    <tr key={props.entityId}>
      <th scope="row">
        <Button
          variant="link"
          onClick={(evt) => nameClick(state, props.entityId)}
        >
          {props.entityName}
        </Button>
      </th>
      <td>{devices}</td>
      <td>
        <span className="btn">
          {props.properties.find((e) => e.id === props.propertyId)?.name}
        </span>
      </td>
      <td>
        <Button
          variant="secondary"
          onClick={(evt) => removeClick(props, props.entityId)}
        >
          Delete
        </Button>
      </td>
    </tr>
  )
}
