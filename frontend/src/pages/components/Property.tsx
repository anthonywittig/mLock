import React from 'react';
import { Button } from 'react-bootstrap';
import { StandardFetch } from '../utils/FetchHelper';

type Adder = (id: string, name: string, updatedBy: string) => void;
type Remover = (id: string) => void;

type Props = {
    entityId: string
    entityName: string,
    updatedBy: string,
    addEntity: Adder,
    removeEntity: Remover,
};

type State = {
    entityId: string,
    entityName: string,
    state: string,
    entityFieldsDisabled: boolean,
};

const Endpoint = "properties";

export class Property extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = this.getResetState();
    }

    getResetState() {
        return {
            entityId: this.props.entityId,
            entityName: this.props.entityName,
            state: this.props.entityId ? "exists" : "new",
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

    newEntitySubmit() {
        if (this.state.entityName === "") {
            // TODO: indicate error.
            return;
        }

        this.setState({
            entityFieldsDisabled: true,
        });

        StandardFetch(Endpoint, {
            method: "POST",
            body: JSON.stringify({
                id: this.state.entityId,
                name: this.state.entityName,
            }),
        })
        .then(response => response.json())
        .then(response => {
            // add to parent
            const e = response.entity;
            this.props.addEntity(e.id, e.name, e.updatedBy);
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
                        <input type="text" className="form-control" id="newEntity" placeholder="Property Name" value={this.state.entityName} onChange={evt => this.updateEntityName(evt)} disabled={this.state.entityFieldsDisabled} onKeyUp={(evt) => evt.key === "Enter" ? this.newEntitySubmit() : ""}/>
                    </th>
                    <td></td>
                    <td><Button variant="secondary" onClick={() => this.newEntitySubmit()} disabled={this.state.entityFieldsDisabled}>Create</Button></td>
                </tr>
            );
        }

        return (
            <tr key={this.props.entityName}>
                <th scope="row">{this.props.entityName}</th>
                <td>{this.props.updatedBy}</td>
                <td><Button variant="secondary" onClick={evt => this.removeClick(this.props.entityId)}>Delete</Button></td>
            </tr>
        );
    }
}