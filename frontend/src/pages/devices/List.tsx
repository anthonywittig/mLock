import React from 'react';
import { Button } from 'react-bootstrap';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';
import { useHistory } from 'react-router';
import { formatDistance } from 'date-fns';

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
    lastRefreshedAt: string,
}

type Property = {
    id: string,
    name: string,
    updatedBy: string,
}

const Endpoint = "devices";

export const List = () => {
    const [entities, setEntities] = React.useState<Entity[]>([]);
    const [loading, setLoading] = React.useState<boolean>(true);
    const [properties, setProperties] = React.useState<Property[]>([]);
    const history = useHistory();

    React.useEffect(() => {
        setLoading(true);

        StandardFetch(Endpoint, {method: "GET"})
        .then(response => response.json())
        .then(response => {
            setEntities(response.entities);
            setLoading(false);
            setProperties(response.extra.properties);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }, []);

    const labelClick = (id: string) => {
        history.push('/' + Endpoint + '/' + id);
    };

    const render = () => {
        return (
            <>
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                        <h2 className="card-title">Details</h2>
                        {renderEntities()}
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
                        <th scope="col">Label</th>
                        <th scope="col">Last Refreshed</th>
                        <th scope="col">Status</th>
                        <th scope="col">Property</th>
                    </tr>
                </thead>
                <tbody>
                    {entities.map(entity =>
                        <tr key={entity.id}>
                            <th scope="row">
                                <Button variant="link" onClick={evt => labelClick(entity.id)}>
                                    {entity.habThing.label}
                                </Button>
                            </th>
                            <td>{ formatDistance(Date.parse(entity.lastRefreshedAt), new Date(), { addSuffix: true }) }</td>
                            <td>{ entity.habThing.statusInfo.status }</td>
                            <td>{ properties.find(e => e.id === entity.propertyId )?.name }</td>
                        </tr>
                    )}
                </tbody>
            </table>
        );
    };

    return render();
};
