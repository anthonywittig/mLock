import React from 'react';
import {
    BrowserRouter as Router,
    Route,
    Switch,
  } from 'react-router-dom';
import { Detail } from './Detail';
import { List } from './List';

export const UnitRoutes = () => {
    return (
        <Router>
            <div>
                <Switch>
                <Route path="/units/:id">
                    <Detail />
                </Route>
                <Route path="/units">
                    <List />
                </Route>
                </Switch>
            </div>
        </Router>
    );
};