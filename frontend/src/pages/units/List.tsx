import React from 'react';
import { Loading } from '../utils/Loading';
import { StandardFetch } from '../utils/FetchHelper';
import { Unit } from './Unit';

type Entity = {
    id: string,
    name: string,
    propertyId: string,
    updatedBy: string,
}

type Property = {
    id: string,
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

    removeEntity(id: string) {
        console.log("removeEntity");
        this.setState({
            entities: this.state.entities.filter(value => {
                return value.id !== id;
            }),
        });
    }

    addEntity(id: string, name: string, propertyId: string, updatedBy: string) {
        this.setState({
            entities: this.state.entities.concat([{
                id,
                name,
                propertyId,
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
                            id={entity.id}
                            entityName={entity.name}
                            propertyId={entity.propertyId}
                            updatedBy={entity.updatedBy}
                            properties={this.state.properties}
                            addEntity={null}
                            removeEntity={id => this.removeEntity(id)}
                        />
                    )}
                    <Unit
                        id=""
                        entityName=""
                        propertyId=""
                        updatedBy=""
                        properties={this.state.properties}
                        addEntity={(id, name, propertyId, updatedBy) => this.addEntity(id, name, propertyId, updatedBy)}
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