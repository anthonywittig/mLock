import React from 'react';
import { Button } from 'react-bootstrap';
import { useHistory } from 'react-router-dom';
import { StandardFetch } from '../utils/FetchHelper';
import { History } from 'history';


type Adder = (name: string, property: string, updatedBy: string) => void;
type NameAction = (name: string) => void;

type Property = {
    name: string,
    createdBy: string,
}

type Props = {
    entityName: string,
    propertyName: string,
    updatedBy: string,
    properties: Property[],
    addEntity: Adder|null,
    removeEntity: NameAction|null,
};

type State = {
    entityFieldsDisabled: boolean,
    setEntityFieldsDisabled: React.Dispatch<React.SetStateAction<boolean>>,
    entityName: string,
    setEntityName: React.Dispatch<React.SetStateAction<string>>,
    propertyName: string,
    setPropertyName: React.Dispatch<React.SetStateAction<string>>,
    entityState: string,
    setEntityState: React.Dispatch<React.SetStateAction<string>>,
    history: History,
};

const Endpoint = "units";

export const Unit = (props: Props) => {
    const state = GetState(props);
    return render(props, state);
};

function GetState(props: Props): State{
    const [entityFieldsDisabled, setEntityFieldsDisabled] = React.useState<boolean>(false);
    const [entityName, setEntityName] = React.useState<string>(props.entityName);
    const [propertyName, setPropertyName] = React.useState<string>(props.propertyName);
    const [entityState, setEntityState] = React.useState<string>(props.entityName ? "exists" : "new");
    const history = useHistory();
    return {
        entityFieldsDisabled,
        setEntityFieldsDisabled,
        entityName,
        setEntityName,
        propertyName,
        setPropertyName,
        entityState,
        setEntityState,
        history,
    };
}

function resetState(props: Props, state: State) {
    // Some dupliation from `getState`...
    state.setEntityFieldsDisabled(false);
    state.setEntityName(props.entityName);
    state.setPropertyName(props.propertyName);
    state.setEntityState(props.entityName ? "exists" : "new");
}

function removeClick(props: Props, name: string) {
    StandardFetch(Endpoint + "/" + encodeURIComponent(name), {method: "DELETE"})
    .then(response => {
        if (response.status === 200) {
            if (props.removeEntity) {
                props.removeEntity(name);
            } else {
                throw new Error("removeEntry is null");
            }
        }
    })
    .catch(err => {
        // TODO: need to indicate error.
        console.log("error: " + err);
    });
}

function nameClick(state: State, name: string) {
   state.history.push('/units/' + encodeURIComponent(name));
}


function updateEntityName(state: State, evt: React.ChangeEvent<HTMLInputElement>) {
    state.setEntityName(evt.target.value);
}


function updatePropertyName(state: State, evt: React.ChangeEvent<HTMLSelectElement>) {
    state.setPropertyName(evt.target.value);
}

function newEntitySubmit(props: Props, state: State) {
    if (!state.entityName || !state.propertyName) {
        // TODO: indicate error.
        return;
    }

    state.setEntityFieldsDisabled(true);

    StandardFetch(Endpoint, {
        method: "POST",
        body: JSON.stringify({ name: state.entityName, propertyName: state.propertyName })
    })
    .then(response => response.json())
    .then(response => {
        // add to parent
        let e = response.entity;
        if (props.addEntity) {
            props.addEntity(e.name, e.propertyName, e.updatedBy);
            resetState(props, state);
        } else {
            throw new Error("addEntity is null");
        }
    })
    .catch(err => {
        // TODO: indicate error.
        state.setEntityFieldsDisabled(false);
    });
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
                        onChange={evt => updateEntityName(state, evt)}
                        disabled={state.entityFieldsDisabled}
                        onKeyUp={(evt) => evt.key === "Enter" ? newEntitySubmit(props, state) : ""}
                    />
                </th>
                <td>
                    <select
                        id="newProperty"
                        className="form-control"
                        onChange={evt => updatePropertyName(state, evt)}
                        disabled={state.entityFieldsDisabled}
                    >
                        <option></option>
                        {props.properties.map(property =>
                            <option value={property.name} selected={property.name === state.propertyName}>
                                {property.name}
                            </option>
                        )}
                    </select>
                </td>
                <td><Button variant="secondary" onClick={() => newEntitySubmit(props, state)} disabled={state.entityFieldsDisabled}>Create</Button></td>
            </tr>
        );
    }

    return (
        <tr key={props.entityName}>
            <th scope="row">
                <Button variant="link" onClick={evt => nameClick(state, props.entityName)}>
                    {props.entityName}
                </Button>
            </th>
            <td>{ props.properties.find(e => e.name === props.propertyName)?.name }</td>
            <td><Button variant="secondary" onClick={evt => removeClick(props, props.entityName)}>Delete</Button></td>
        </tr>
    );
}
