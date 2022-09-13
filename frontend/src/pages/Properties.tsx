import React from 'react'
import { Loading } from './utils/Loading'
import { StandardFetch } from './utils/FetchHelper'
import { Property } from './components/Property'

type Entity = {
    id: string
    name: string
    updatedBy: string
}

type Props = {}

type State = {
    entities: Entity[]
    loadingEntities: boolean
}

export class Properties extends React.Component<Props, State> {
    state: Readonly<State> = {
        entities: [],
        loadingEntities: true,
    }

    componentDidMount() {
        StandardFetch("properties", {method: "GET"})
        .then(response => response.json())
        .then(response => {
            this.setState({
                loadingEntities: false,
                entities: response.entities,
            })
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err)
        })
    }

    removeEntity(id: string) {
        this.setState({
            entities: this.state.entities.filter(value => {
                return value.id !== id
            }),
        })
    }

    addEntity(id: string, name: string, updatedBy: string) {
        this.setState({
            entities: this.state.entities.concat([{
                id,
                name,
                updatedBy,
            }]),
        })
    }

    renderEntitiesTable() {
        if (this.state.loadingEntities) {
            return <Loading />
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
                        <Property entityId={entity.id} entityName={entity.name} updatedBy={entity.updatedBy} removeEntity={id => this.removeEntity(id)} addEntity={props => console.log("should never happen")}/>
                    )}
                    <Property entityId="" entityName="" updatedBy="" removeEntity={id => console.log("should never happen")} addEntity={(id, name, updatedBy) => this.addEntity(id, name, updatedBy)}/>
                </tbody>
            </table>
        )
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
        )
    }
}