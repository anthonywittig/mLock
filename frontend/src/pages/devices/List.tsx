import React from 'react';
import { Badge, Button, OverlayTrigger, Tooltip } from 'react-bootstrap';
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
                        <th scope="col">Status</th>
                        <th scope="col">Battery</th>
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
                            <td>{ renderEntityStatus(entity) }</td>
                            <td>{ renderEntityBatteryLevel(entity) }</td>
                            <td>{ properties.find(e => e.id === entity.propertyId )?.name }</td>
                            <td>{ units.find(e => e.id === entity.unitId )?.name }</td>
                            <td>{ renderDeleteButton(entity) }</td>
                        </tr>
                    )}
                </tbody>
            </table>
        );
    };

    const renderDeleteButton = (entity: DeviceT) => {
        const lr = Date.parse(entity.lastRefreshedAt);
        const recently = sub(new Date(), {hours: 2});

        if (isAfter(lr, recently) && entity.habThing.statusInfo.status === "ONLINE") {
            return (
                <OverlayTrigger overlay={<Tooltip id="tooltip-disabled">Must be offline or not refreshed in the last two hours.</Tooltip>}>
                    <span className="d-inline-block">
                        <Button variant="secondary" disabled style={{ pointerEvents: 'none' }}>Delete</Button>
                    </span>
                </OverlayTrigger>
            );
        }

        return <Button variant="secondary" onClick={() => deleteDevice(entity.id)}>Delete</Button>;
    };

    const renderEntityStatus = (entity: DeviceT) => {
        const warnings = getLastRefreshedWarnings(entity);
        warnings.push.apply(warnings, getStatusInfoWarning(entity));
        warnings.push.apply(warnings, getLastWentOfflineWarnings(entity));

        return (
            <ul>
                { warnings.map(warn =>
                    <li>{ warn }</li>
                )}
            </ul>
        );
    };

    const renderEntityBatteryLevel = (entity : DeviceT) => {

        const lu = entity.battery.lastUpdatedAt;
        if (lu === null) {
            return <></>;
        }

        const lud = Date.parse(lu);
        const recently = sub(new Date(), {days: 1, hours: 12});

        if (isBefore(lud, recently)) {
            const distance = formatDistance(lud, new Date(), { addSuffix: true });
            return <Badge variant="danger">Last battery level taken over { distance }</Badge>;
        }

        const level = parseFloat(entity.battery.level);
        if (isNaN(level)) {
            return <Badge variant="danger">Unknown</Badge>;
        }

        if (level < 25) {
            return <Badge variant="danger">{ level }%</Badge>;
        }

        return <>{ level }%</>;
    };

    const getLastRefreshedWarnings = (entity : DeviceT) => {
        const warnings: JSX.Element[] = [];

        const lr = Date.parse(entity.lastRefreshedAt);
        const recently = sub(new Date(), {minutes: 10});

        if (isBefore(lr, recently)) {
            const distance = formatDistance(lr, new Date(), { addSuffix: true });
            warnings.push(<Badge variant="danger">Last Data Sync: { distance }</Badge>);
        }

        return warnings;
    };

    const getLastWentOfflineWarnings = (entity: DeviceT) => {
        const warnings: JSX.Element[] = [];

        const lwo = entity.lastWentOfflineAt;
        if (lwo === null) {
            return warnings;
        }

        const lwod = Date.parse(entity.lastWentOfflineAt!);
        const recently = sub(new Date(), {days: 1});

        if (entity.habThing.statusInfo.status !== "ONLINE" || isAfter(lwod, recently)) {
            const distance = formatDistance(lwod, new Date(), { addSuffix: true });
            warnings.push(<Badge variant="danger">Was Offline: { distance }</Badge>);
        }

        return warnings;
    };

    const getStatusInfoWarning = (entity: DeviceT) => {
        const warnings: JSX.Element[] = [];

        const s = entity.habThing.statusInfo.status;
        if (s !== "ONLINE") {
            warnings.push(<Badge variant="danger">{ s }</Badge>);
        }

        const sd = entity.habThing.statusInfo.statusDetail;
        if (sd !== "NONE") {
            warnings.push(<Badge variant="secondary">{ sd }</Badge>);
        }

        return warnings;
    };

    return render();
};
