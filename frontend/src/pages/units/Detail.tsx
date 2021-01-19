import React from 'react';
import { Button, Form} from 'react-bootstrap';
import { useRouteMatch } from 'react-router-dom';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';

type Entity = {
    id: string,
    name: string,
    propertyId: string,
    calendarUrl: string,
    updatedBy: string,
}

type Property = {
    id: string,
    name: string,
    createdBy: string,
}

type State = {
    entity: Entity,
    setEntity: React.Dispatch<React.SetStateAction<Entity>>,
    loading: boolean,
    setLoading: React.Dispatch<React.SetStateAction<boolean>>,
    properties: Property[],
    setProperties: React.Dispatch<React.SetStateAction<Property[]>>,
};

type MatchParams = {id: string};

const Endpoint = "units";

export const Detail = () => {
    const state = GetState();
    const m = useRouteMatch('/units/:id');
    const mp = m?.params as MatchParams;
    console.log(mp.id);

    const didMountRef = React.useRef(true);
    React.useEffect(() => {
        if (didMountRef.current) {
            xComponentDidMount(state, mp.id);
        }
        didMountRef.current = false;
    }, [state, mp.id]);

    return render(state);
};

function GetState(): State{
    const [entity, setEntity] = React.useState<Entity>({
        id: "",
        name: "",
        propertyId: "",
        calendarUrl: "",
        updatedBy: "",
    });
    const [loading, setLoading] = React.useState<boolean>(true);
    const [properties, setProperties] = React.useState<Property[]>([]);
    console.log("Get state");
    console.log(entity);
    return {
        entity, setEntity,
        loading, setLoading,
        properties, setProperties,
    };
}

function xComponentDidMount(state: State, entityId: string) {
    state.setLoading(true);

    StandardFetch(Endpoint + "/" + entityId, {method: "GET"})
    .then(response => response.json())
    .then(response => {
        state.setEntity(response.entity);
        state.setLoading(false);
        state.setProperties(response.extra.properties);
    })
    .catch(err => {
        // TODO: indicate error.
        console.log(err);
    });
}

function detailFormNameChange(state: State, evt: React.ChangeEvent<HTMLInputElement>) {
    let entity = copyEntity(state.entity);
    entity.name = evt.target.value;
    state.setEntity(entity);
}

function detailFormPropertyIdChange(state: State, evt: React.ChangeEvent<HTMLSelectElement>) {
    let entity = copyEntity(state.entity);
    entity.propertyId = evt.target.value;
    state.setEntity(entity);
}

function detailFormCalendarUrlChange(state: State, evt: React.ChangeEvent<HTMLSelectElement>) {
    let entity = copyEntity(state.entity);
    entity.calendarUrl = evt.target.value;
    state.setEntity(entity);
}

function copyEntity(old: Entity): Entity{
    return {
        id: old.id,
        name: old.name,
        propertyId: old.propertyId,
        calendarUrl: old.calendarUrl,
        updatedBy: old.updatedBy,
    };
}

function detailFormSubmit(state: State, evt: React.FormEvent<HTMLFormElement>) {
    evt.preventDefault();

    state.setLoading(true);

    StandardFetch(Endpoint + "/" + state.entity.id, {
        method: "PUT",
        body: JSON.stringify({
            name: state.entity.name,
            propertyId: state.entity.propertyId,
            calendarUrl: state.entity.calendarUrl,
        })
    })
    .then(response => response.json())
    .then(response => {
        state.setEntity(response.entity);
        state.setLoading(false);
    })
    .catch(err => {
        // TODO: indicate error.
        console.log(err);
    });
}

function render(state: State) {
    return (
        <div>
            <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                <div className="card-body">
                <h2 className="card-title">Details</h2>
                {renderEntity(state)}
                </div>
            </div>
        </div>
    );
}

function renderEntity(state: State) {
    if (state.loading) {
        return <Loading />;
    }
    return (
        <Form onSubmit={evt => detailFormSubmit(state, evt)}>
            <Form.Group>
                <Form.Label>Name</Form.Label>
                <Form.Control type="text" value={state.entity.name} onChange={evt => detailFormNameChange(state, evt as any)}/>
            </Form.Group>

            <Form.Group controlId="exampleForm.ControlSelect1">
                <Form.Label>Property</Form.Label>
                <Form.Control as="select" onChange={evt => detailFormPropertyIdChange(state, evt as any)}>
                    {state.properties.map(property =>
                        <option value={property.id} selected={property.id === state.entity.propertyId}>
                            {property.name}
                        </option>
                    )}
                </Form.Control>
            </Form.Group>

            <Form.Group>
                <Form.Label>Calendar URL</Form.Label>
                <Form.Control type="text" value={state.entity.calendarUrl} onChange={evt => detailFormCalendarUrlChange(state, evt as any)}/>
            </Form.Group>

            <Button variant="secondary" type="submit">
                Submit
            </Button>
        </Form>
    );
}
