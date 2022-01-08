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
    }, [entities.length]);

    const deleteDevice = (id: string) => {
        setLoading(true);

        StandardFetch(Endpoint + "/" + id, {
            method: "DELETE",
        })
        .then(_ => {
            setEntities(
                entities.filter(value => {
                    return value.id !== id;
                }),
            );
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
                        <h2 className="card-title">Devices</h2>
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
                                    { entity.rawDevice.name }
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
        const recently = sub(new Date(), {minutes: 20});

        if (isAfter(lr, recently)) {
            return (
                <OverlayTrigger overlay={<Tooltip id="tooltip-disabled">The device was recently pulled from the controller.</Tooltip>}>
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
        if (!entity.rawDevice.battery.batteryPowered) {
            return <></>;
        }

        const lu = entity.lastRefreshedAt;
        const lud = Date.parse(lu);
        const recently = sub(new Date(), {days: 1, hours: 12});
        const level = entity.rawDevice.battery.level;

        if (isBefore(lud, recently) || level === null) {
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
        const recently = sub(new Date(), {minutes: 20});

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

        const recently = sub(new Date(), {days: 1});
        const lwond = Date.parse(entity.lastWentOnlineAt!);
        const lwoffd = Date.parse(entity.lastWentOfflineAt!);

        if (entity.rawDevice.status !== "ONLINE") {
            const distance = formatDistance(lwoffd, new Date(), { addSuffix: true });
            warnings.push(<Badge variant="danger">Went Offline: { distance }</Badge>);
        } else if (isAfter(lwond, recently)) {
            const distance = formatDistance(lwond, new Date(), { addSuffix: true });
            warnings.push(<Badge variant="danger">Went Online: { distance }</Badge>);
        }

        return warnings;
    };

    const getStatusInfoWarning = (entity: DeviceT) => {
        const warnings: JSX.Element[] = [];

        const s = entity.rawDevice.status;
        if (s !== "ONLINE") {
            warnings.push(<Badge variant="danger">{ s[0].toUpperCase() + s.slice(1) }</Badge>);
        }

        return warnings;
    };

    return render();
};
