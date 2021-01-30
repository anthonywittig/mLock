import React from 'react';
import { Button, Form} from 'react-bootstrap';
import { useRouteMatch } from 'react-router-dom';
import { format, parseISO } from 'date-fns';
import { useHistory } from 'react-router-dom';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';

type Entity = {
    name: string,
    propertyName: string,
    calendarUrl: string,
    updatedBy: string,
}

type Reservation = {
    id: string,
    start: string,
    startDate: Date,
    end: string,
    endDate: Date,
    summary: string,
    status: string,
}

type Property = {
    name: string,
    createdBy: string,
}

type MatchParams = {name: string};

const Endpoint = "units";

export const Detail = () => {
    const [entity, setEntity] = React.useState<Entity>({
        name: "",
        propertyName: "",
        calendarUrl: "",
        updatedBy: "",
    });
    const [loading, setLoading] = React.useState<boolean>(true);
    const [properties, setProperties] = React.useState<Property[]>([]);
    const [reservations, setReservations] = React.useState<Reservation[]>([]);
    const history = useHistory();

    const m = useRouteMatch('/units/:name');
    const mp = m?.params as MatchParams;
    const name = mp.name;

    React.useEffect(() => {
        setLoading(true);

        StandardFetch(Endpoint + "/" + encodeURIComponent(name), {method: "GET"})
        .then(response => response.json())
        .then(response => {
            setEntity(response.entity);
            setLoading(false);
            setProperties(response.extra.properties);

            let reservations = response.extra.reservations as Reservation[];
            reservations.forEach(r => {
                // The dates are naive, so cut off the zone.
                r.startDate = parseISO(r.start.slice(0, -1));
                r.endDate = parseISO(r.end.slice(0, -1));
            });
            setReservations(reservations);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }, [name]);

    const detailFormNameChange = (evt: React.ChangeEvent<HTMLInputElement>) => {
        setEntity({
            ...entity,
            name: evt.target.value,
        });
    };

    const detailFormPropertyNameChange = (evt: React.ChangeEvent<HTMLSelectElement>) => {
        setEntity({
            ...entity,
            propertyName: evt.target.value,
        });
    };

    const detailFormCalendarUrlChange = (evt: React.ChangeEvent<HTMLSelectElement>) => {
        setEntity({
            ...entity,
            calendarUrl: evt.target.value,
        });
    };

    const detailFormSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
        evt.preventDefault();

        setLoading(true);

        StandardFetch(Endpoint + "/" + encodeURIComponent(name), {
            method: "PUT",
            /*body: JSON.stringify({
                name: entity.name,
                propertyName: entity.propertyName,
                calendarUrl: entity.calendarUrl,
            })*/
            body: JSON.stringify(entity)
        })
        .then(response => response.json())
        .then(response => {
            setEntity(response.entity);
            setLoading(false);
            history.push('/units/' + encodeURIComponent(response.entity.name));
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
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                        <h2 className="card-title">Upcoming Reservations</h2>
                        {renderCalendar()}
                    </div>
                </div>
            </>
        );
    };

    const renderCalendar = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <table className="table table-responsive-sm">
                <thead>
                    <tr>
                        <th scope="col">Transaction #</th>
                        <th scope="col">Start Date</th>
                        <th scope="col">End Date</th>
                    </tr>
                </thead>
                <tbody>
                    {reservations.map(res =>
                        <tr>
                            <th scope="row">{res.summary}</th>
                            <td>{format(res.startDate, "LL/dd/yyyy")}</td>
                            <td>{format(res.endDate, "LL/dd/yyyy")}</td>
                        </tr>
                    )}
                </tbody>
            </table>
        );
    };

    const renderEntity = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <Form onSubmit={evt => detailFormSubmit(evt)}>
                <Form.Group>
                    <Form.Label>Name</Form.Label>
                    <Form.Control type="text" value={entity.name} onChange={evt => detailFormNameChange(evt as any)}/>
                </Form.Group>

                <Form.Group controlId="exampleForm.ControlSelect1">
                    <Form.Label>Property</Form.Label>
                    <Form.Control as="select" onChange={evt => detailFormPropertyNameChange(evt as any)}>
                        {properties.map(property =>
                            <option value={property.name} selected={property.name === entity.propertyName}>
                                {property.name}
                            </option>
                        )}
                    </Form.Control>
                </Form.Group>

                <Form.Group>
                    <Form.Label>Calendar URL</Form.Label>
                    <Form.Control type="text" value={entity.calendarUrl} onChange={evt => detailFormCalendarUrlChange(evt as any)}/>
                </Form.Group>

                <Button variant="secondary" type="submit">
                    Submit
                </Button>
            </Form>
        );
    };

    return render();
};
