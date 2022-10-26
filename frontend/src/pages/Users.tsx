import React from "react"
import { Button, Table } from "react-bootstrap"
import { Loading } from "./utils/Loading"
import { StandardFetch } from "./utils/FetchHelper"

type User = {
  id: string
  email: string
  updatedBy: string
}

type Props = {}

type State = {
  users: User[]
  newUser: string
  newUserFieldEnabled: boolean
  newUserButtonEnabled: boolean
  loadingUsers: boolean
}

export class Users extends React.Component<Props, State> {
  state: Readonly<State> = {
    users: [],
    newUser: "",
    newUserFieldEnabled: true,
    newUserButtonEnabled: false,
    loadingUsers: true,
  }

  componentDidMount() {
    StandardFetch("users", { method: "GET" })
      .then((response) => response.json())
      .then((response) => {
        this.setState({
          loadingUsers: false,
          users: response.Users,
        })
      })
      .catch((err) => {
        // Might want to retry once on failure.
        console.log(err)
      })
  }

  removeUserClick(id: string) {
    this.setState({ loadingUsers: true })

    StandardFetch("users/" + id, { method: "DELETE" })
      .then((response) => response.json())
      .then((response) => {
        if (response.Users) {
          this.setState({ users: response.Users })
        }
        this.setState({ loadingUsers: false })
      })
      .catch((err) => {
        // Need to indicate error...
        console.log("error: " + err)
      })
  }

  newUserClick() {
    this.setState({
      newUserFieldEnabled: false,
      newUserButtonEnabled: false,
    })

    StandardFetch("users", {
      method: "POST",
      body: JSON.stringify({ email: this.state.newUser }),
    })
      .then((response) => response.json())
      .then((response) => {
        this.setState({
          users: response.Users,
          newUser: "",
          newUserFieldEnabled: true,
        })
      })
      .catch((err) => {
        // Need to indicate error...
        this.setState({
          newUserFieldEnabled: true,
          newUserButtonEnabled: true, // Not that helpful but probably less confusing?
        })
      })
  }

  updateNewUserValue(evt: React.ChangeEvent<HTMLInputElement>) {
    this.setState({
      newUser: evt.target.value,
      newUserButtonEnabled: evt.target.value !== "",
    })
  }

  renderUsersTable() {
    if (this.state.loadingUsers) {
      return <Loading />
    }
    return (
      <Table className="table table-responsive-sm">
        <thead>
          <tr>
            <th scope="col">Email Address</th>
            <th scope="col">Created By</th>
            <th scope="col">Actions</th>
          </tr>
        </thead>
        <tbody>
          {this.state.users.map((user) => (
            <tr key={user.id}>
              <th scope="row">{user.email}</th>
              <td>{user.updatedBy}</td>
              <td>
                <Button
                  variant="secondary"
                  onClick={(evt) => this.removeUserClick(user.id)}
                >
                  Delete
                </Button>
              </td>
            </tr>
          ))}
          <tr key="newUser">
            <th scope="row">
              <input
                type="text"
                className="form-control"
                id="newUser"
                placeholder="Enter new user's Google email address"
                value={this.state.newUser}
                onChange={(evt) => this.updateNewUserValue(evt)}
                disabled={!this.state.newUserFieldEnabled}
                onKeyUp={(evt) =>
                  evt.key === "Enter" ? this.newUserClick() : ""
                }
              />
            </th>
            <td></td>
            <td>
              <Button
                variant="secondary"
                onClick={() => this.newUserClick()}
                disabled={!this.state.newUserButtonEnabled}
              >
                Create
              </Button>
            </td>
          </tr>
        </tbody>
      </Table>
    )
  }

  render() {
    return (
      <div>
        <div
          className="card"
          style={{ marginBottom: "1rem", marginTop: "1rem" }}
        >
          <div className="card-body">
            <h2 className="card-title">Users</h2>
            {this.renderUsersTable()}
          </div>
        </div>
      </div>
    )
  }
}
