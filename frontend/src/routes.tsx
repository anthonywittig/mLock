import {
  BrowserRouter as Router,
  Redirect,
  Route,
  Switch,
} from "react-router-dom"
import { DeviceRoutes } from "./pages/devices/Routes"
import { Navigation } from "./navigation"
import { Properties } from "./pages/Properties"
import { PrivacyPolicy } from "./pages/PrivacyPolicy"
import { TermsOfService } from "./pages/TermsOfService"
import { SignIn } from "./pages/SignIn"
import { UnitRoutes } from "./pages/units/Routes"
import { Users } from "./pages/Users"

export const Routes = () => {
  return (
    <Router>
      <Navigation />
      <div>
        <Switch>
          <Route path="/devices">
            <DeviceRoutes />
          </Route>
          <Route path="/properties">
            <Properties />
          </Route>
          <Route path="/privacy-policy">
            <PrivacyPolicy />
          </Route>
          <Route path="/sign-in">
            <SignIn />
          </Route>
          <Route path="/terms-of-service">
            <TermsOfService />
          </Route>
          <Route path="/units">
            <UnitRoutes />
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
  )
}
