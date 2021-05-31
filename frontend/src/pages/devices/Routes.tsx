import React from 'react';
import {
    BrowserRouter as Router,
    Route,
    Switch,
  } from 'react-router-dom';
import { Detail } from './Detail';
import { List } from './List';

const Endpoint = "devices";

export const DeviceRoutes = () => {
    return (
        <Router>
            <div>
                <Switch>
                <Route path={"/" + Endpoint + "/:id"}>
                    <Detail />
                </Route>
                <Route path={"/" + Endpoint}>
                    <List />
                </Route>
                </Switch>
            </div>
        </Router>
    );
};