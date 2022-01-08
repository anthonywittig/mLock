import {
    Route,
    Switch,
  } from 'react-router-dom';
import { Detail } from './Detail';
import { List } from './List';

export const UnitRoutes = () => {
    return (
        <Switch>
        <Route path="/units/:id">
            <Detail />
        </Route>
        <Route path="/units">
            <List />
        </Route>
        </Switch>
    );
};