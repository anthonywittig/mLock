import React from 'react'
import {
    BrowserRouter as Router,
    Redirect,
    Route,
    Switch,
  } from 'react-router-dom';
import { Home } from './pages/Home'
import { PrivacyPolicy } from './pages/PrivacyPolicy'
import { TermsOfService } from './pages/TermsOfService'
import { SignIn } from './pages/SignIn'
import { Users } from './pages/Users'

export const Routes = () => {
    return (
        <Router>
            <div>
                {/* A <Switch> looks through its children <Route>s and
                    renders the first one that matches the current URL. */}
                <Switch>
                <Route path="/about">
                    <About />
                </Route>
                <Route path="/privacy-policy">
                    <PrivacyPolicy />
                </Route>
                <Route path="/sign-in">
                    <SignIn />
                </Route>
                <Route path="/terms-of-service">
                    <TermsOfService/>
                </Route>
                <Route path="/users">
                    <Users />
                </Route>
                <Route path="/">
                    <Redirect to="/sign-in" />
                </Route>
                </Switch>
            </div>
        </Router>
    );
}

function About() {
    return <h2>About</h2>;
}