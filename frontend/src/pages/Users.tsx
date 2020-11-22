import React from 'react'
import { Button } from 'react-bootstrap';

type Props = {
    users: string[];
};

type State = {
    users: string[];
};

export class Users extends React.Component<Props, State> {
    static defaultProps: Props = {
        users: []
    }

    state: Readonly<State> = {
        users: []
    }

    componentDidMount() {
        fetch("https://api.zcclock.com/users", {
            "method": "GET",
            "headers": {
                "apikey": "apikey",
            }
        })
        .then(response => response.json())
        .then(response => {
            console.log(response);
            this.setState({
                users: response.Users
            });
        })
        .catch(err => {
            console.log(err); 
        });
    }

    render() {
        return (
            <div>
                <div className="card" style={{marginBottom: "1rem", marginTop: "1rem"}}>
                    <div className="card-body">
                    <h2 className="card-title">Users</h2>
                    <table className="table table-responsive-sm">
                        <thead>
                        <tr>
                            <th scope="col">Email Address</th>
                            <th scope="col">Actions</th>
                        </tr>
                        </thead>
                        <tbody>
                            {this.state.users.map(user =>
                                <tr>
                                    <th scope="row">---{user}----</th>
                                    <td><Button variant="secondary">Remove</Button></td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                    <Button variant="secondary">Add User</Button>
                    </div>
                </div>
            </div>
        );
    }
}