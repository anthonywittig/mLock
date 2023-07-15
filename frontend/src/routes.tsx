import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom"
import { DeviceRoutes } from "./pages/devices/Routes"
import { Navigation } from "./navigation"
import { Properties } from "./pages/Properties"
import { PrivacyPolicy } from "./pages/PrivacyPolicy"
import { TermsOfService } from "./pages/TermsOfService"
import { SignIn } from "./pages/SignIn"
import { UnitRoutes } from "./pages/units/Routes"
import { Users } from "./pages/Users"

export const NavRoutes = () => {
  return (
    <Router>
      <Navigation />
      <div>
        <Routes>
          <Route path="/devices/*" element={<DeviceRoutes />} />
          <Route path="/properties/*" element={<Properties />} />
          <Route path="/privacy-policy/*" element={<PrivacyPolicy />} />
          <Route path="/sign-in/*" element={<SignIn />} />\
          <Route path="/terms-of-service/*" element={<TermsOfService />} />
          <Route path="/units/*" element={<UnitRoutes />} />
          <Route path="/users/*" element={<Users />} />
          <Route path="/" element={<Navigate to="/devices/" replace />} />
        </Routes>
      </div>
    </Router>
  )
}
