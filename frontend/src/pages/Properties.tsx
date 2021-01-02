import React from 'react';
import { Button } from 'react-bootstrap';
import { Loading } from './utils/Loading';
import { StandardFetch } from './utils/FetchHelper';

type Entity = {
    ID: string;
    Name: string;
    CreatedBy: string;
}

type Props = {};

type State = {
    entities: Entity[];
    newEntity: string;
    newEntityFieldEnabled: boolean;
    newEntityButtonEnabled: boolean;
    loadingEntities: boolean;
};

export class Properties extends React.Component<Props, State> {
    state: Readonly<State> = {
        entities: [],
        newEntity: "",
        newEntityFieldEnabled: true,
        newEntityButtonEnabled: false,
        loadingEntities: true,
    }

    componentDidMount() {
        StandardFetch("properties", {method: "GET"})
        .then(response => response.json())
        .then(response => {
            this.setState({
                loadingEntities: false,
                entities: response.Entities,
            });
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }

    removeEntityClick(id: string) {
        this.setState({loadingEntities: true});

        StandardFetch("properties/" + id, {method: "DELETE"})
        .then(response => response.json())
        .then(response => {
            if (response.Entities) {
                this.setState({
                    entities: response.Entities,
                });
            }
            this.setState({
                loadingEntities: false,
            });
        })
        .catch(err => {
            // TODO: need to indicate error.
            console.log("error: " + err);
        });
    }

    newEntityClick() {
        this.setState({
            newEntityFieldEnabled: false,
            newEntityButtonEnabled: false,
        });

        StandardFetch("properties", {
            method: "POST",
            body: JSON.stringify({ name: this.state.newEntity })
        })
        .then(response => response.json())
        .then(response => {
            this.setState({
                entities: response.Entities,
                newEntity: "",
                newEntityFieldEnabled: true,
            });
        })
        .catch(err => {
            // TODO: indicate error.
            this.setState({
                newEntityFieldEnabled: true,
                newEntityButtonEnabled: true,
            });
        });
    }

    updateNewEntityValue(evt: React.ChangeEvent<HTMLInputElement>) {
        this.setState({
            newEntity: evt.target.value,
            newEntityButtonEnabled: evt.target.value !== "",
        });
    }

    renderEntitiesTable() {
        if (this.state.loadingEntities) {
            return <Loading />;
        }
        return (
            <table className="table table-responsive-sm">
                <thead>
                    <tr>
                        <th scope="col">Name</th>
                        <th scope="col">Created By</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {this.state.entities.map(entity =>
                        <tr key={entity.Name}>
                            <th scope="row">{entity.Name}</th>
                            <td>{entity.CreatedBy}</td>
                            <td><Button variant="secondary" onClick={evt => this.removeEntityClick(entity.ID)}>Delete</Button></td>
                        </tr>
                    )}
                    <tr key="newEntity">
                        <th scope="row">
                            <input type="text" className="form-control" id="newEntity" placeholder="Property Name" value={this.state.newEntity} onChange={evt => this.updateNewEntityValue(evt)} disabled={!this.state.newEntityFieldEnabled} onKeyUp={(evt) => evt.key === "Enter" ? this.newEntityClick() : ""}/>
                        </th>
                        <td></td>
                        <td><Button variant="secondary" onClick={() => this.newEntityClick()} disabled={!this.state.newEntityButtonEnabled}>Create</Button></td>
                    </tr>
                </tbody>
            </table>
        );
    }


    render() {
        return (
            <div>
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                    <h2 className="card-title">Properties</h2>
                    {this.renderEntitiesTable()}
                    </div>
                </div>
            </div>
        );
    }
}