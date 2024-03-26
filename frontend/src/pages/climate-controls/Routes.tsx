import { Routes, Route } from "react-router-dom"
import { Detail } from "./Detail"
import { List } from "./List"

export const ClimateControlRoutes = () => {
  return (
    <Routes>
      <Route path={":id"} element={<Detail />} />
      <Route path={""} element={<List />} />
    </Routes>
  )
}
