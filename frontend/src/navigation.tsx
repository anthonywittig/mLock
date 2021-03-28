import React from 'react';
import { Navbar, Nav } from 'react-bootstrap';
import { useLocation } from 'react-router-dom';

export const Navigation = () => {
    let location = useLocation();

    if (location.pathname.startsWith("/sign-in")) {
        return null;
    }

    return (
        <Navbar bg="light" expand="lg">
            <Navbar.Brand href="/units/">MLock</Navbar.Brand>
            <Navbar.Toggle aria-controls="basic-navbar-nav" />
            <Navbar.Collapse id="basic-navbar-nav">
                <Nav className="mr-auto">
                    <Nav.Link href="/units/" className={location.pathname.startsWith("/units/") ? "active" : ""}>Units</Nav.Link>
                    <Nav.Link href="/properties/" className={location.pathname.startsWith("/properties/") ? "active" : ""}>Properties</Nav.Link>
                    <Nav.Link href="/users/" className={location.pathname.startsWith("/users/") ? "active" : ""}>Users</Nav.Link>
                </Nav>
            </Navbar.Collapse>
        </Navbar>
    );
};