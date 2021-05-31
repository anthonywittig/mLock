import React from 'react';
import { Form} from 'react-bootstrap';
import { useRouteMatch } from 'react-router-dom';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';

type Entity = {
    id: string,
    propertyId: string,
    habThing: {
        label: string,
        statusInfo: {
            status: string,
            statusDetail: string,
        },
        uid: string,
    },
    lastRefreshedAt: Date,
}

type Property = {
    id: string,
    name: string,
    updatedBy: string,
}

type MatchParams = {id: string};

const Endpoint = "devices";

export const Detail = () => {
    const [entity, setEntity] = React.useState<Entity>({
        id: "",
        propertyId: "",
        habThing: {
            label: "",
            statusInfo: {
                status: "",
                statusDetail: "",
            },
            uid: "",
        },
        lastRefreshedAt: new Date(),
    });
    const [loading, setLoading] = React.useState<boolean>(true);
    const [properties, setProperties] = React.useState<Property[]>([]);

    const m = useRouteMatch('/' + Endpoint + '/:id');
    const mp = m?.params as MatchParams;
    const id = mp.id;

    React.useEffect(() => {
        setLoading(true);

        StandardFetch(Endpoint + "/" + id, {method: "GET"})
        .then(response => response.json())
        .then(response => {
            setEntity(response.entity);
            setLoading(false);
            setProperties(response.extra.properties);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }, [id]);


    const render = () => {
        return (
            <>
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                        <h2 className="card-title">Details</h2>
                        {renderEntity()}
                    </div>
                </div>
            </>
        );
    };

    const renderEntity = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <Form>
                <Form.Group>
                    <Form.Label>Label</Form.Label>
                    <Form.Control type="text" value={entity.habThing.label} disabled={true} />
                </Form.Group>

                <Form.Group>
                    <Form.Label>Last Refreshed At</Form.Label>
                    <Form.Control type="text" value={entity.lastRefreshedAt.toString()} disabled={true}/>
                </Form.Group>

                <Form.Group>
                    <Form.Label>Status</Form.Label>
                    <Form.Control type="text" value={entity.habThing.statusInfo.status} disabled={true}/>
                </Form.Group>

                <Form.Group controlId="property">
                    <Form.Label>Property</Form.Label>
                    <Form.Control as="select" disabled={true}>
                        {properties.map(property =>
                            <option value={property.id} selected={property.id === entity.propertyId}>
                                {property.name}
                            </option>
                        )}
                    </Form.Control>
                </Form.Group>
            </Form>
        );
    };

    return render();
};
