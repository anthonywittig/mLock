import React from 'react';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';
import { Unit } from './Unit';

type Entity = {
    name: string,
    propertyName: string,
    updatedBy: string,
}

type Property = {
    name: string,
    createdBy: string,
}

type Props = {};

type State = {
    entities: Entity[],
    loadingEntities: boolean,
    properties: Property[],
};

const Endpoint = "units";

export class List extends React.Component<Props, State> {
    state: Readonly<State> = {
        entities: [],
        loadingEntities: true,
        properties: [],
    }

    componentDidMount() {
        StandardFetch(Endpoint, {method: "GET"})
        .then(response => response.json())
        .then(response => {
            this.setState({
                entities: response.entities,
                loadingEntities: false,
                properties: response.extra.properties,
            });
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    }

    removeEntity(name: string) {
        this.setState({
            entities: this.state.entities.filter(value => {
                return value.name !== name;
            }),
        });
    }

    addEntity(name: string, propertyName: string, updatedBy: string) {
        this.setState({
            entities: this.state.entities.concat([{
                name,
                propertyName,
                updatedBy,
            }]),
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
                        <th scope="col">Property</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {this.state.entities.map(entity =>
                        <Unit
                            entityName={entity.name}
                            propertyName={entity.propertyName}
                            updatedBy={entity.updatedBy}
                            properties={this.state.properties}
                            addEntity={null}
                            removeEntity={name => this.removeEntity(name)}
                        />
                    )}
                    <Unit
                        entityName=""
                        propertyName=""
                        updatedBy=""
                        properties={this.state.properties}
                        addEntity={(name, propertyName, updatedBy) => this.addEntity(name, propertyName, updatedBy)}
                        removeEntity={null}
                    />
                </tbody>
            </table>
        );
    }

    render() {
        return (
            <div>
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                    <h2 className="card-title">Units</h2>
                    {this.renderEntitiesTable()}
                    </div>
                </div>
            </div>
        );
    }
}