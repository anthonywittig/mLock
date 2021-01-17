import React from 'react';
import { Button } from 'react-bootstrap';
import { StandardFetch } from '../utils/FetchHelper';
import { Redirect, RouteComponentProps } from "react-router-dom";


type Adder = (id: string, name: string, property: string, updatedBy: string) => void;
type IdAction = (id: string) => void;

type Property = {
    id: string,
    name: string,
    createdBy: string,
}

type Props = RouteComponentProps & {
    id: string,
    entityName: string,
    propertyId: string,
    updatedBy: string,
    properties: Property[],
    addEntity: Adder|null,
    removeEntity: IdAction|null,
};

type State = {
    entityFieldsDisabled: boolean,
    entityName: string,
    propertyId: string,
    redirect: string,
    state: string,
};

const Endpoint = "units";

export class Unit extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = this.getResetState();
    }

    getResetState(): State {
        return {
            entityFieldsDisabled: false,
            entityName: this.props.entityName,
            propertyId: this.props.propertyId,
            redirect: "",
            state: this.props.id ? "exists" : "new",
        };
    }

    removeClick(id: string) {
        StandardFetch(Endpoint + "/" + id, {method: "DELETE"})
        .then(response => {
            if (response.status === 200) {
                if (this.props.removeEntity) {
                    this.props.removeEntity(id);
                } else {
                    throw new Error("removeEntry is null");
                }
            }
        })
        .catch(err => {
            // TODO: need to indicate error.
            console.log("error: " + err);
        });
    }

    nameClick(id: string) {
        /*
        if (this.props.navToDetail) {
            this.props.navToDetail(id);
        }
        */
       console.log("nameClick - fix me!");
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
            if (this.props.addEntity) {
                this.props.addEntity(e.id, e.name, e.propertyId, e.updatedBy);
                this.setState(this.getResetState());
            } else {
                throw new Error("addEntity is null");
            }
        })
        .catch(err => {
            // TODO: indicate error.
            this.setState({
                entityFieldsDisabled: false,
            });
        });
    }

    render() {
        if (this.state.redirect) {
            return <Redirect to={this.state.redirect} />;
        }
        
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
                    <td><Button variant="secondary" onClick={() => this.newEntitySubmit()} disabled={this.state.entityFieldsDisabled}>Create</Button></td>
                </tr>
            );
        }

        return (
            <tr key={this.props.id}>
                <th scope="row">
                    <Button variant="link" onClick={evt => this.nameClick(this.props.id)}>
                        {this.props.entityName}
                    </Button>
                </th>
                <td>{ this.props.properties.find(e => e.id === this.props.propertyId)?.name }</td>
                <td><Button variant="secondary" onClick={evt => this.removeClick(this.props.id)}>Delete</Button></td>
            </tr>
        );
    }
}