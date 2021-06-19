import React from 'react';
import { Form} from 'react-bootstrap';
import { formatDistance } from 'date-fns';
import { useRouteMatch } from 'react-router-dom';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';

type Entity = {
    id: string,
    propertyId: string,
    habThing: {
        configuration: {
            usercode_code_1: string,
            usercode_code_10: string,
            usercode_code_11: string,
            usercode_code_12: string,
            usercode_code_13: string,
            usercode_code_14: string,
            usercode_code_15: string,
            usercode_code_16: string,
            usercode_code_17: string,
            usercode_code_18: string,
            usercode_code_19: string,
            usercode_code_2: string,
            usercode_code_20: string,
            usercode_code_21: string,
            usercode_code_22: string,
            usercode_code_23: string,
            usercode_code_24: string,
            usercode_code_25: string,
            usercode_code_26: string,
            usercode_code_27: string,
            usercode_code_28: string,
            usercode_code_29: string,
            usercode_code_3: string,
            usercode_code_30: string,
            usercode_code_4: string,
            usercode_code_5: string,
            usercode_code_6: string,
            usercode_code_7: string,
            usercode_code_8: string,
            usercode_code_9: string,
        },
        label: string,
        statusInfo: {
            status: string,
            statusDetail: string,
        },
        uid: string,
    },
    lastRefreshedAt: string,
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
            configuration: {
                usercode_code_1: "",
                usercode_code_10: "",
                usercode_code_11: "",
                usercode_code_12: "",
                usercode_code_13: "",
                usercode_code_14: "",
                usercode_code_15: "",
                usercode_code_16: "",
                usercode_code_17: "",
                usercode_code_18: "",
                usercode_code_19: "",
                usercode_code_2: "",
                usercode_code_20: "",
                usercode_code_21: "",
                usercode_code_22: "",
                usercode_code_23: "",
                usercode_code_24: "",
                usercode_code_25: "",
                usercode_code_26: "",
                usercode_code_27: "",
                usercode_code_28: "",
                usercode_code_29: "",
                usercode_code_3: "",
                usercode_code_30: "",
                usercode_code_4: "",
                usercode_code_5: "",
                usercode_code_6: "",
                usercode_code_7: "",
                usercode_code_8: "",
                usercode_code_9: "",
            },
            label: "",
            statusInfo: {
                status: "",
                statusDetail: "",
            },
            uid: "",
        },
        lastRefreshedAt: "",
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
                    <Form.Label>Last Refreshed</Form.Label>
                    <Form.Control type="text" value={formatDistance(Date.parse(entity.lastRefreshedAt), new Date(), { addSuffix: true }) } disabled={true}/>
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

                {
                    // Don't look :(
                    [
                        ["code 1", entity.habThing.configuration.usercode_code_1],
                        ["code 2", entity.habThing.configuration.usercode_code_2],
                        ["code 3", entity.habThing.configuration.usercode_code_3],
                        ["code 4", entity.habThing.configuration.usercode_code_4],
                        ["code 5", entity.habThing.configuration.usercode_code_5],
                        ["code 6", entity.habThing.configuration.usercode_code_6],
                        ["code 7", entity.habThing.configuration.usercode_code_7],
                        ["code 8", entity.habThing.configuration.usercode_code_8],
                        ["code 9", entity.habThing.configuration.usercode_code_9],
                        ["code 10", entity.habThing.configuration.usercode_code_10],
                        ["code 11", entity.habThing.configuration.usercode_code_11],
                        ["code 12", entity.habThing.configuration.usercode_code_12],
                        ["code 13", entity.habThing.configuration.usercode_code_13],
                        ["code 14", entity.habThing.configuration.usercode_code_14],
                        ["code 15", entity.habThing.configuration.usercode_code_15],
                        ["code 16", entity.habThing.configuration.usercode_code_16],
                        ["code 17", entity.habThing.configuration.usercode_code_17],
                        ["code 18", entity.habThing.configuration.usercode_code_18],
                        ["code 19", entity.habThing.configuration.usercode_code_19],
                        ["code 20", entity.habThing.configuration.usercode_code_20],
                        ["code 21", entity.habThing.configuration.usercode_code_21],
                        ["code 22", entity.habThing.configuration.usercode_code_22],
                        ["code 23", entity.habThing.configuration.usercode_code_23],
                        ["code 24", entity.habThing.configuration.usercode_code_24],
                        ["code 25", entity.habThing.configuration.usercode_code_25],
                        ["code 26", entity.habThing.configuration.usercode_code_26],
                        ["code 27", entity.habThing.configuration.usercode_code_27],
                        ["code 28", entity.habThing.configuration.usercode_code_28],
                        ["code 29", entity.habThing.configuration.usercode_code_29],
                        ["code 30", entity.habThing.configuration.usercode_code_30],
                    ].map((u) => {
                        return (
                            <Form.Group style={u[1] ? {} : {display: "none"}}>
                                <Form.Label>{u[0]}</Form.Label>
                                <Form.Control type="text" value={u[1]} disabled={true}/>
                            </Form.Group>
                        );
                    })
                }
            </Form>
        );
    };

    return render();
};
