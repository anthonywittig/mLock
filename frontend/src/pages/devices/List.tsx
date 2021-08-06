import React from 'react';
import { Badge, Button } from 'react-bootstrap';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';
import { useHistory } from 'react-router';
import { formatDistance, isAfter, isBefore, sub } from 'date-fns';

type Property = {
    id: string,
    name: string,
    updatedBy: string,
}

const Endpoint = "devices";

export const List = () => {
    const [entities, setEntities] = React.useState<DeviceT[]>([]);
    const [loading, setLoading] = React.useState<boolean>(true);
    const [properties, setProperties] = React.useState<Property[]>([]);
    const [units, setUnits] = React.useState<UnitT[]>([]);
    const history = useHistory();

    React.useEffect(() => {
        setLoading(true);

        StandardFetch(Endpoint, {method: "GET"})
        .then(response => response.json())
        .then(response => {
            setEntities(response.entities);
            setLoading(false);
            setProperties(response.extra.properties);
            setUnits(response.extra.units);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }, []);

    const deleteDevice = (id: string) => {
        setLoading(true);

        StandardFetch(Endpoint + "/" + id, {
            method: "DELETE",
        })
        .then(_ => {
            history.push('/' + Endpoint + '/');
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    };

    const labelClick = (id: string) => {
        history.push('/' + Endpoint + '/' + id);
    };

    const render = () => {
        return (
            <>
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                        <h2 className="card-title">Details</h2>
                        { renderEntities() }
                    </div>
                </div>
            </>
        );
    };

    const renderEntities = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <table className="table table-responsive-sm">
                <thead>
                    <tr>
                        <th scope="col">Name</th>
                        <th scope="col">Last Refreshed</th>
                        <th scope="col">Status</th>
                        <th scope="col">Went Offline</th>
                        <th scope="col">Property</th>
                        <th scope="col">Unit</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {entities.map(entity =>
                        <tr key={ entity.id }>
                            <th scope="row">
                                <Button variant="link" onClick={evt => labelClick(entity.id)}>
                                    { entity.habThing.label }
                                </Button>
                            </th>
                            <td>{ renderEntityLastRefreshed(entity) }</td>
                            <td>{ renderEntityStatus(entity) }</td>
                            <td>{ renderEntityLastWentOffline(entity) }</td>
                            <td>{ properties.find(e => e.id === entity.propertyId )?.name }</td>
                            <td>{ units.find(e => e.id === entity.unitId )?.name }</td>
                            <td><Button variant="secondary" onClick={() => deleteDevice(entity.id)}>Delete</Button></td>
                        </tr>
                    )}
                </tbody>
            </table>
        );
    };

    const renderEntityLastRefreshed = (entity : DeviceT) => {
        const lr = Date.parse(entity.lastRefreshedAt);
        const recently = sub(new Date(), {minutes: 10});
        const distance = formatDistance(lr, new Date(), { addSuffix: true });

        if (isBefore(lr, recently)) {
            return <Badge variant="danger">{ distance }</Badge>;
        }

        return "recently";
    };

    const renderEntityStatus = (entity : DeviceT) => {
        const status = entity.habThing.statusInfo.status;
        if (status !== "ONLINE") {
            return <Badge variant="danger">{ entity.habThing.statusInfo.status }</Badge>;
        }
        return entity.habThing.statusInfo.status;
    };

    const renderEntityLastWentOffline = (entity : DeviceT) => {
        const lwo = entity.lastWentOfflineAt;

        if (lwo === null) {
            return "";
        }

        const lwod = Date.parse(entity.lastWentOfflineAt!);
        const recently = sub(new Date(), {days: 1});
        const distance = formatDistance(lwod, new Date(), { addSuffix: true });

        if (isAfter(lwod, recently)) {
            return <Badge variant="danger">{ distance }</Badge>;
        }

        return distance;
    };

    return render();
};
