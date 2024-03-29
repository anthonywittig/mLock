import React from "react"
import { Navbar, Nav, Button } from "react-bootstrap"
import { useLocation } from "react-router-dom"
import { googleLogout } from "@react-oauth/google"

export const Navigation = () => {
  let location = useLocation()

  if (location.pathname.startsWith("/sign-in")) {
    return null
  }

  return (
    <Navbar bg="light" expand="lg">
      <Navbar.Brand href="/devices/">MLock</Navbar.Brand>
      <Navbar.Toggle aria-controls="basic-navbar-nav" />
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="mr-auto">
          <Nav.Link
            href="/units/"
            className={location.pathname.startsWith("/units/") ? "active" : ""}
          >
            Units
          </Nav.Link>
          <Nav.Link
            href="/properties/"
            className={
              location.pathname.startsWith("/properties/") ? "active" : ""
            }
          >
            Properties
          </Nav.Link>
          <Nav.Link
            href="/devices/"
            className={
              location.pathname.startsWith("/devices/") ? "active" : ""
            }
          >
            Devices
          </Nav.Link>
          <Nav.Link
            href="/climate-controls/"
            className={
              location.pathname.startsWith("/climate-controls/") ? "active" : ""
            }
          >
            Climate Controls
          </Nav.Link>
          <Nav.Link
            href="/users/"
            className={location.pathname.startsWith("/users/") ? "active" : ""}
          >
            Users
          </Nav.Link>
        </Nav>
      </Navbar.Collapse>
      <Nav>
        <Button
          variant="outline-secondary"
          onClick={googleLogout}
          href="/sign-in/"
          className="justify-content-md-end"
        >
          Log Out
        </Button>
      </Nav>
    </Navbar>
  )
}
