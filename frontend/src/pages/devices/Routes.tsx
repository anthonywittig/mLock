import {
    Route,
    Switch,
  } from 'react-router-dom'
import { Detail } from './Detail'
import { List } from './List'

const Endpoint = "devices"

export const DeviceRoutes = () => {
    return (
        <Switch>
        <Route path={"/" + Endpoint + "/:id"}>
            <Detail />
        </Route>
        <Route path={"/" + Endpoint}>
            <List />
        </Route>
        </Switch>
    )
}