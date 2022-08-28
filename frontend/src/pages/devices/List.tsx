import React from 'react';
import { Badge, Button, OverlayTrigger, Tooltip } from 'react-bootstrap';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';
import { useHistory } from 'react-router';
import { formatDistance, isAfter, isBefore, sub } from 'date-fns';

const Endpoint = "devices";

export const List = () => {
    const [entities, setEntities] = React.useState<DeviceT[]>([]);
    const [loading, setLoading] = React.useState<boolean>(true);
    const [units, setUnits] = React.useState<UnitT[]>([]);
    const history = useHistory();

    React.useEffect(() => {
        setLoading(true);

        StandardFetch(Endpoint, {method: "GET"})
        .then(response => response.json())
        .then(response => {
            setEntities(response.entities);
            setLoading(false);
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
                        <th scope="col">Online</th>
                        <th scope="col">Status</th>
                        <th scope="col">Battery</th>
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
                            <td>{ renderOnline(entity) }</td>
                            <td>{ renderEntityStatus(entity) }</td>
                            <td>{ renderEntityBatteryLevel(entity) }</td>
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
        // This should really be something much smaller, like 20 minutes, but since we have periods of time where we don't sync for an hour, we need something at least 60 minutes long.
        const recently = sub(new Date(), {minutes: 130});

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
        warnings.push.apply(warnings, getLockResponsivenessWarnings(entity));

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

    const renderOnline = (entity: DeviceT) => {
        if (entity.rawDevice.status === "ONLINE") {
            return <p>Online</p>;
        }else if (entity.rawDevice.status !== "ONLINE") {
            return <Badge variant="danger">Offline</Badge>;
        }else{
            return <Badge variant="warning">Error!</Badge>;
        }
    };

    const getLastRefreshedWarnings = (entity : DeviceT) => {
        const warnings: JSX.Element[] = [];

        const lr = Date.parse(entity.lastRefreshedAt);
        const recently = sub(new Date(), {minutes: 70});

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

    const getLockResponsivenessWarnings = (entity: DeviceT) => {
        const warnings: JSX.Element[] = [];

        if (entity.rawDevice.status !== "ONLINE") {
            return warnings;
        }

        const tooSoon = sub(new Date(), {minutes: 10});
        const expectedResponseInMinutes = 20;

        for (let i = 0; i < entity.managedLockCodes.length; i++) {
            const lc = entity.managedLockCodes[i];

            // We really should consider the `startedRemovingAt` and `wasCompletedAt` timestamps, but right now we're only syncing every hour during the time that most codes are being removed.
            if (lc.startedAddingAt) {
                const sa = Date.parse(lc.startedAddingAt);
                if (isBefore(sa, tooSoon)) {
                    if (warnings.length && lc.status !== "Adding") {
                        // Once we have one warning, we'll only add additional ones for `Adding` codes.
                        continue;
                    }
                    if (lc.wasEnabledAt) {
                        const wc = Date.parse(lc.wasEnabledAt);
                        const minutesBetween = (wc - sa) / 1000 / 60;
                        if (expectedResponseInMinutes < minutesBetween) {
                            const distance = formatDistance(sa, wc);
                            warnings.push(<Badge variant="danger">Slow to Respond (took { distance } to add code { lc.code })</Badge>);
                        }
                    } else {
                        warnings.push(<Badge variant="danger">Not Responding (for code { lc.code })</Badge>);
                    }
                }
            }
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
