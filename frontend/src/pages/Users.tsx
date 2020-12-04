import React from 'react'
import { Button } from 'react-bootstrap';
import { Loading } from './utils/Loading'

type User = {
    ID: string;
    Email: string;
}

type Props = {};

type State = {
    users: User[];
    newUser: string;
    newUserFieldEnabled: boolean;
    newUserButtonEnabled: boolean;
    loadingUsers: boolean;
};

export class Users extends React.Component<Props, State> {
    /*
    static defaultProps: Props = {
        users: [],
        newUser: "",
        newUserFieldEnabled: true,
        newUserButtonEnabled: false,
    }
    */

    state: Readonly<State> = {
        users: [],
        newUser: "",
        newUserFieldEnabled: true,
        newUserButtonEnabled: false,
        loadingUsers: true,
    }

    constructor(props: Props) {
        super(props);
        /*
        this.state = {
            users: props.users,
            newUser: props.newUser,
            newUserFieldEnabled: props.newUserFieldEnabled,
            newUserButtonEnabled: props.newUserButtonEnabled,
        };
        */
      }

    componentDidMount() {
        fetch("https://api2.zcclock.com/users", {
            "method": "GET",
            "headers": {
                "apikey": "apikey",
            }
        })
        .then(response => response.json())
        .then(response => {
            console.log(response);
            this.setState({
                loadingUsers: false,
                users: response.Users
            });
        })
        .catch(err => {
            // Might want to retry once on failure.
            console.log(err);
        });
    }

    newUserClick() {
        this.setState({
            newUserFieldEnabled: false,
            newUserButtonEnabled: false,
        });
        console.log(this.state.newUser);

        fetch("https://api2.zcclock.com/users", {
            method: "POST",
            headers: {
                "apikey": "apikey",
            },
            body: JSON.stringify({ email: this.state.newUser })
        })
        .then(response => response.json())
        .then(response => {
            console.log(response);
            this.setState({
                users: response.Users,
                newUser: "",
                newUserFieldEnabled: true,
            });
        })
        .catch(err => {
            // Need to indicate error...
            this.setState({
                newUserFieldEnabled: true,
                newUserButtonEnabled: true, // Not that helpful but probably less confusing?
            });
        });
    }

    updateNewUserValue(evt: React.ChangeEvent<HTMLInputElement>) {
        this.setState({
            newUser: evt.target.value,
            newUserButtonEnabled: evt.target.value !== "",
        })
    }

    renderUsersTable() {
        if (this.state.loadingUsers) {
            return <Loading />;
        }
        return (
            <table className="table table-responsive-sm">
                <thead>
                    <tr>
                        <th scope="col">Email Address</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {this.state.users.map(user =>
                        <tr key={user.Email}>
                            <th scope="row">{user.Email}</th>
                            <td><Button variant="secondary">Remove</Button></td>
                        </tr>
                    )}
                    <tr key="newUser">
                        <th scope="row">
                            <input type="text" className="form-control" id="newUser" placeholder="Enter Google email address" value={this.state.newUser} onChange={evt => this.updateNewUserValue(evt)} disabled={!this.state.newUserFieldEnabled} />
                        </th>
                        <td><Button variant="secondary" onClick={() => this.newUserClick()} disabled={!this.state.newUserButtonEnabled}>Add User</Button></td>
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
                    <h2 className="card-title">Users</h2>
                    {this.renderUsersTable()}
                    </div>
                </div>
            </div>
        );
    }
}