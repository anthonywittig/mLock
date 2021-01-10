import React from 'react';
import { Detail } from './Detail';
import { List } from './List';

type Props = {};
type State = {
    path: string,
};

export class UnitsRouter extends React.Component<Props, State> {
    state: Readonly<State> = {
        path: "",
    }

    constructor(props: Props) {
        super(props);

        const path = window.location.pathname.replace("/units/", "");
        this.state = {
            path,
        };
    }

    updatePath(path: string) {
        this.setState({
            path,
        });
    }

    render() {
        if (this.state.path === "") {
            return <List updatePath={p => this.updatePath(p)}/>;
        }
        return <Detail entityId={this.state.path} />;
    }
}