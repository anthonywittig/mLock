import React from 'react';
import { Button, Form} from 'react-bootstrap';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';

type Entity = {
    id: string,
    name: string,
    propertyId: string,
    calendarUrl: string,
    updatedBy: string,
}

type Property = {
    id: string,
    name: string,
    createdBy: string,
}

type Props = {
    entityId: string,
};

type State = {
    entity: Entity,
    loading: boolean,
    properties: Property[],
};

const Endpoint = "units";

export class Detail extends React.Component<Props, State> {
    state: Readonly<State> = {
        entity: {id: "", name: "", propertyId: "", calendarUrl: "", updatedBy: ""},
        loading: true,
        properties: [],
    }

    componentDidMount() {
        StandardFetch(Endpoint + "/" + this.props.entityId, {method: "GET"})
        .then(response => response.json())
        .then(response => {
            this.setState({
                entity: response.entity,
                loading: false,
                properties: response.extra.properties,
            });
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }
    
    detailFormNameChange(evt: React.ChangeEvent<HTMLInputElement>) {
        let entity = this.state.entity;
        entity.name = evt.target.value;
        this.setState({entity});
    }

    detailFormPropertyIdChange(evt: React.ChangeEvent<HTMLSelectElement>) {
        let entity = this.state.entity;
        entity.propertyId = evt.target.value;
        this.setState({entity});
    }

    detailFormCalendarUrlChange(evt: React.ChangeEvent<HTMLSelectElement>) {
        let entity = this.state.entity;
        entity.calendarUrl = evt.target.value;
        this.setState({entity});
    }

    detailFormSubmit(evt: React.FormEvent<HTMLFormElement>) {
        evt.preventDefault();

        this.setState({loading: true});

        StandardFetch(Endpoint + "/" + this.props.entityId, {
            method: "PUT",
            body: JSON.stringify({
                name: this.state.entity.name,
                propertyId: this.state.entity.propertyId,
                calendarUrl: this.state.entity.calendarUrl,
            })
        })
        .then(response => response.json())
        .then(response => {
            this.setState({
                entity: response.entity,
                loading: false,
            });
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }

    render() {
        return (
            <div>
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                    <h2 className="card-title">Details</h2>
                    {this.renderEntity()}
                    </div>
                </div>
            </div>
        );
    }

    renderEntity() {
        if (this.state.loading) {
            return <Loading />;
        }
        return (
            <Form onSubmit={evt => this.detailFormSubmit(evt)}>
                <Form.Group>
                    <Form.Label>Name</Form.Label>
                    <Form.Control type="text" value={this.state.entity.name} onChange={evt => this.detailFormNameChange(evt as any)}/>
                </Form.Group>

                <Form.Group controlId="exampleForm.ControlSelect1">
                    <Form.Label>Property</Form.Label>
                    <Form.Control as="select" onChange={evt => this.detailFormPropertyIdChange(evt as any)}>
                        {this.state.properties.map(property =>
                            <option value={property.id} selected={property.id === this.state.entity.propertyId}>
                                {property.name}
                            </option>
                        )}
                    </Form.Control>
                </Form.Group>

                <Form.Group>
                    <Form.Label>Calendar URL</Form.Label>
                    <Form.Control type="text" value={this.state.entity.calendarUrl} onChange={evt => this.detailFormCalendarUrlChange(evt as any)}/>
                </Form.Group>

                <Button variant="secondary" type="submit">
                    Submit
                </Button>
            </Form>
        );
    }
}