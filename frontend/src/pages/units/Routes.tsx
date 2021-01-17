import React from 'react';
import {
    BrowserRouter as Router,
    Route,
    RouteComponentProps,
    Switch,
  } from 'react-router-dom';
import { Detail } from './Detail';
import { List } from './List';

type Props = RouteComponentProps; 

export const UnitRoutes = (props: Props) => {
    return (
        <Router>
            <div>
                <Switch>
                <Route path="/units/:id">
                    <Detail
                        history={props.history}
                        location={props.location}
                        match={props.match}
                    />
                </Route>
                <Route path="/units">
                    <List
                        history={props.history}
                        location={props.location}
                        match={props.match}
                    />
                </Route>
                </Switch>
            </div>
        </Router>
    );
};