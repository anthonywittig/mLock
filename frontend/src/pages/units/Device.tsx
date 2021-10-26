import React from 'react';
import { Button } from 'react-bootstrap';
import { useHistory } from 'react-router-dom';
import { StandardFetch } from '../utils/FetchHelper';
import { History } from 'history';

type IdAction = (id: string) => void;

type Props = {
    id: string,
    updatedBy: string,
    devices: DeviceT[],
    addEntity: IdAction|null,
    removeEntity: IdAction|null,
};

type State = {
    entityFieldsDisabled: boolean,
    setEntityFieldsDisabled: React.Dispatch<React.SetStateAction<boolean>>,
    entityId: string,
    setEntityId: React.Dispatch<React.SetStateAction<string>>,
    history: History,
};

const Endpoint = "device";

export const Device = (props: Props) => {
    const state = GetState(props);
    return render(props, state);
};

function GetState(props: Props): State{
    const [entityFieldsDisabled, setEntityFieldsDisabled] = React.useState<boolean>(false);
    const [entityId, setEntityId] = React.useState<string>(props.id);
    const history = useHistory();

    return {
        entityFieldsDisabled,
        setEntityFieldsDisabled,
        entityId,
        setEntityId,
        history,
    };
}

function resetState(props: Props, state: State) {
    // Some duplication from `getState`...
    state.setEntityFieldsDisabled(false);
    state.setEntityId(props.id);
}

function updateEntityId(state: State, evt: React.ChangeEvent<HTMLSelectElement>) {
    state.setEntityId(evt.target.value);
}

function addSubmit(props: Props, state: State) {
    if (!state.entityId) {
        // TODO: indicate error.
        return;
    }

    state.setEntityFieldsDisabled(true);

    if (!props.addEntity) {
        throw new Error("addEntity is null");
    }
    props.addEntity(props.id);

    resetState(props, state);
}

function labelClick(state: State, id: string) {
    state.history.push('/' + Endpoint + '/' + id);
}

function removeClick(props: Props, id: string) {
    if (!props.removeEntity) {
        throw new Error("removeEntry is null");
    }
    props.removeEntity(id);
}

function render(props: Props, state: State) {
    if (!state.entityId) {
        return (
            <tr key="newEntity">
                <th scope="row">
                    <select
                        id="newDevice"
                        className="form-control"
                        onChange={evt => updateEntityId(state, evt)}
                        disabled={state.entityFieldsDisabled}
                    >
                        <option></option>
                        {props.devices.map(device =>
                            <option value={device.id} selected={device.id === state.entityId}>
                                {device.rawDevice.name}
                            </option>
                        )}
                    </select>
                </th>
                <td><Button variant="secondary" onClick={() => addSubmit(props, state)} disabled={state.entityFieldsDisabled}>Add</Button></td>
            </tr>
        );
    }

    return (
        <tr key={props.id}>
            <th scope="row">
                <Button variant="link" onClick={evt => labelClick(state, props.id)}>
                    { props.devices.find(e => e.id === props.id)?.rawDevice.name}
                </Button>
            </th>
            <td><Button variant="secondary" onClick={evt => removeClick(props, props.id)}>Remove</Button></td>
        </tr>
    );
}
