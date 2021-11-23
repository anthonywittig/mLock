import React from 'react';
import { Button, Form} from 'react-bootstrap';
import { formatDistance } from 'date-fns';
import { useHistory, useRouteMatch } from 'react-router-dom';
import { LockCode } from './components/LockCode';
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
        lastRefreshedAt: "",
        lastWentOfflineAt: null,
        lastWentOnlineAt: null,
        managedLockCodes: [],
        rawDevice: {
            battery: {
                batteryPowered: false,
                level: 0,
            },
            categoryId: "",
            lockCodes: null,
            name: "",
            status: "",
        }
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
                    <div className="card-body">
                        <h2 className="card-title">Current Lock Codes</h2>
                        {renderCurrentLockCodes()}
                    </div>
                    <div className="card-body">
                        <h2 className="card-title">Add Lock Code</h2>
                        <LockCode deviceId={entity.id} managedLockCode={null}/>
                    </div>
                </div>
            </>
        );
    };

    const renderCurrentLockCodes = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <>
                {
                    entity.managedLockCodes.map((lc) => {
                        return (
                            <div><LockCode deviceId={entity.id} managedLockCode={lc}/><br/></div>
                        );
                    })
                }
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
                    <Form.Control type="text" value={entity.rawDevice.name} disabled={true} />
                </Form.Group>

                <Form.Group>
                    <Form.Label>Last Refreshed</Form.Label>
                    <Form.Control type="text" value={formatDistance(Date.parse(entity.lastRefreshedAt), new Date(), { addSuffix: true }) } disabled={true}/>
                </Form.Group>

                <Form.Group>
                    <Form.Label>Status</Form.Label>
                    <Form.Control type="text" value={entity.rawDevice.status} disabled={true}/>
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

                <Button variant="secondary" type="submit">
                    Update
                </Button>
            </Form>
        );
    };

    return render();
};
