import React from 'react';
import { Loading } from './utils/Loading';
import { StandardFetch } from './utils/FetchHelper';
import { Property } from './components/Property';

type Entity = {
    ID: string;
    Name: string;
    CreatedBy: string;
}

type Props = {};

type State = {
    entities: Entity[];
    newEntity: string;
    loadingEntities: boolean;
};

export class Properties extends React.Component<Props, State> {
    state: Readonly<State> = {
        entities: [],
        newEntity: "",
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

    removeEntity(id: string) {
        this.setState({
            entities: this.state.entities.filter(value => {
                return value.ID !== id;
            }),
        });
    }

    addEntity(id: string, name: string, createdBy: string) {
        this.setState({
            entities: this.state.entities.concat([{
                ID: id,
                Name: name,
                CreatedBy: createdBy,
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
                        <th scope="col">Last Updated By</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {this.state.entities.map(entity =>
                        <Property id={entity.ID} entityName={entity.Name} createdBy={entity.CreatedBy} removeEntity={id => this.removeEntity(id)} addEntity={props => console.log("should never happen")}/>
                    )}
                    <Property id="" entityName="" createdBy="" removeEntity={id => console.log("should never happen")} addEntity={(id, name, createdBy) => this.addEntity(id, name, createdBy)}/>
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