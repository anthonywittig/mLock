import React from 'react';
import { Loading } from './utils/Loading';
import { StandardFetch } from './utils/FetchHelper';
import { Unit } from './components/Unit';

type Entity = {
    id: string,
    name: string,
    propertyId: string,
    createdBy: string,
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

export class Units extends React.Component<Props, State> {
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
        this.setState({
            entities: this.state.entities.filter(value => {
                return value.id !== id;
            }),
        });
    }

    addEntity(id: string, name: string, propertyId: string, createdBy: string) {
        this.setState({
            entities: this.state.entities.concat([{
                id,
                name,
                propertyId,
                createdBy,
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
                        <th scope="col">Last Updated By</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {this.state.entities.map(entity =>
                        <Unit
                            id={entity.id}
                            entityName={entity.name}
                            propertyId={entity.propertyId}
                            createdBy={entity.createdBy}
                            properties={this.state.properties}
                            addEntity={props => console.log("should never happen")}
                            removeEntity={id => this.removeEntity(id)}
                        />
                    )}
                    <Unit
                        id=""
                        entityName=""
                        propertyId=""
                        createdBy=""
                        properties={this.state.properties}
                        addEntity={(id, name, propertyId, createdBy) => this.addEntity(id, name, propertyId, createdBy)}
                        removeEntity={id => console.log("should never happen")}
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