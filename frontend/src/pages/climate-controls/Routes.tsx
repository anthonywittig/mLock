import { Route, Routes } from "react-router-dom"
import { List } from "./List"

export const ClimateControlRoutes = () => {
  return (
    <Routes>
      <Route path={""} element={<List />} />
    </Routes>
  )
}
