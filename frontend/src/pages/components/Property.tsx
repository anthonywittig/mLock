import React from 'react';
import { Button } from 'react-bootstrap';
import { StandardFetch } from '../utils/FetchHelper';

type Adder = (name: string, createdBy: string) => void;
type Remover = (name: string) => void;

type Props = {
    entityName: string,
    createdBy: string,
    addEntity: Adder,
    removeEntity: Remover,
};

type State = {
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
            entityName: this.props.entityName,
            state: this.props.entityName ? "exists" : "new",
            entityFieldsDisabled: false,
        };
    }

    removeClick(name: string) {
        StandardFetch(Endpoint + "/" + encodeURIComponent(name), {method: "DELETE"})
        .then(response => {
            if (response.status === 200) {
                this.props.removeEntity(name);
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
            body: JSON.stringify({ name: this.state.entityName })
        })
        .then(response => response.json())
        .then(response => {
            // add to parent
            const e = response.entity;
            this.props.addEntity(e.name, e.createdBy);
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
                <td>{this.props.createdBy}</td>
                <td><Button variant="secondary" onClick={evt => this.removeClick(this.props.entityName)}>Delete</Button></td>
            </tr>
        );
    }
}