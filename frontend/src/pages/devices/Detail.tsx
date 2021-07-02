import React from 'react';
import { Button, Form} from 'react-bootstrap';
import { formatDistance } from 'date-fns';
import { useHistory, useRouteMatch } from 'react-router-dom';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';

type Property = {
    id: string,
    name: string,
    updatedBy: string,
}

type MatchParams = {id: string};

const Endpoint = "devices";

export const Detail = () => {
    const [entity, setEntity] = React.useState<DeviceT>({
        id: "",
        propertyId: "",
        unitId: "",
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
    const [units, setUnits] = React.useState<UnitT[]>([]);

    const m = useRouteMatch('/' + Endpoint + '/:id');
    const mp = m?.params as MatchParams;
    const id = mp.id;
    const history = useHistory();

    React.useEffect(() => {
        setLoading(true);

        StandardFetch(Endpoint + "/" + id, {method: "GET"})
        .then(response => response.json())
        .then(response => {
            setEntity(response.entity);
            setLoading(false);
            setProperties(response.extra.properties);
            setUnits(response.extra.units);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }, [id]);

    const detailFormUnitChange = (evt: React.ChangeEvent<HTMLSelectElement>) => {
        let val: (string | null) = evt.target.value;
        if (val === "") {
            val = null;
        }
        setEntity({
            ...entity,
            unitId: val,
        });
    };

    const formSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
        evt.preventDefault();

        setLoading(true);

        StandardFetch(Endpoint + "/" + id, {
            method: "PUT",
            body: JSON.stringify(entity)
        })
        .then(response => response.json())
        .then(response => {
            setEntity(response.entity);
            setLoading(false);
            history.push('/' + Endpoint + '/' + response.entity.id);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    };

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
            <Form onSubmit={evt => formSubmit(evt)}>
                <Form.Group>
                    <Form.Label>Name</Form.Label>
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

                <Form.Group controlId="unit">
                    <Form.Label>Unit</Form.Label>
                    <Form.Control as="select" onChange={evt => detailFormUnitChange(evt as any)}>
                        <option></option>
                        {units.filter(unit =>
                            unit.propertyId === entity.propertyId
                        ).map(unit =>
                            <option value={unit.id} selected={unit.id === entity.unitId}>
                                {unit.name}
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

                <Button variant="secondary" type="submit">
                    Update
                </Button>

            </Form>
        );
    };

    return render();
};
