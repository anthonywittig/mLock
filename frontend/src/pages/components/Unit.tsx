import React from 'react';
import { Button } from 'react-bootstrap';
import { StandardFetch } from '../utils/FetchHelper';

type Adder = (id: string, name: string, property: string, updatedBy: string) => void;
type Remover = (id: string) => void;

type Property = {
    id: string,
    name: string,
    createdBy: string,
}

type Props = {
    id: string,
    entityName: string,
    propertyId: string,
    createdBy: string,
    properties: Property[],
    addEntity: Adder,
    removeEntity: Remover,
};

type State = {
    entityName: string,
    propertyId: string,
    state: string,
    entityFieldsDisabled: boolean,
};

const Endpoint = "units";

export class Unit extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = this.getResetState();
    }

    getResetState(): State {
        return {
            entityName: this.props.entityName,
            propertyId: this.props.propertyId,
            state: this.props.id ? "exists" : "new",
            entityFieldsDisabled: false,
        };
    }

    removeClick(id: string) {
        StandardFetch(Endpoint + "/" + id, {method: "DELETE"})
        .then(response => {
            if (response.status === 200) {
                this.props.removeEntity(id);
            }
        })
        .catch(err => {
            // TODO: need to indicate error.
            console.log("error: " + err);
        });
    }

    updateEntityName(evt: React.ChangeEvent<HTMLInputElement>) {
        this.setState({
            entityName: evt.target.value,
        });
    }

    updatePropertyId(evt: React.ChangeEvent<HTMLSelectElement>) {
        this.setState({
            propertyId: evt.target.value,
        });
    }

    newEntitySubmit() {
        if (!this.state.entityName || !this.state.propertyId) {
            // TODO: indicate error.
            return;
        }

        this.setState({
            entityFieldsDisabled: true,
        });

        StandardFetch(Endpoint, {
            method: "POST",
            body: JSON.stringify({ name: this.state.entityName, propertyId: this.state.propertyId })
        })
        .then(response => response.json())
        .then(response => {
            // add to parent
            let e = response.entity;
            this.props.addEntity(e.id, e.name, e.propertyId, e.createdBy);
            this.setState(this.getResetState());
        })
        .catch(err => {
            // TODO: indicate error.
            this.setState({
                entityFieldsDisabled: false,
            });
        });
    }

    render() {
        if (this.state.state === "new") {
            return (
                <tr key="newEntity">
                    <th scope="row">
                        <input
                            type="text"
                            className="form-control"
                            id="newName"
                            placeholder="Name"
                            value={this.state.entityName}
                            onChange={evt => this.updateEntityName(evt)}
                            disabled={this.state.entityFieldsDisabled}
                            onKeyUp={(evt) => evt.key === "Enter" ? this.newEntitySubmit() : ""}
                        />
                    </th>
                    <td>
                        <select
                            id="newProperty"
                            className="form-control"
                            onChange={evt => this.updatePropertyId(evt)}
                            disabled={this.state.entityFieldsDisabled}
                        >
                            <option></option>
                            {this.props.properties.map(property =>
                                <option value={property.id} selected={property.id === this.state.propertyId}>
                                    {property.name}
                                </option>
                            )}
                        </select>
                    </td>
                    <td></td>
                    <td><Button variant="secondary" onClick={() => this.newEntitySubmit()} disabled={this.state.entityFieldsDisabled}>Create</Button></td>
                </tr>
            );
        }

        return (
            <tr key={this.props.id}>
                <th scope="row">{this.props.entityName}</th>
                <td>{ this.props.properties.find(e => e.id === this.props.propertyId)?.name }</td>
                <td>{this.props.createdBy}</td>
                <td><Button variant="secondary" onClick={evt => this.removeClick(this.props.id)}>Delete</Button></td>
            </tr>
        );
    }
}