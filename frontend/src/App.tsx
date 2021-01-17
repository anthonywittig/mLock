import React from 'react';
import './App.css';
import { RouteComponentProps } from "react-router-dom";
import { Routes } from './routes';

type Props = RouteComponentProps;

function App(props: Props) {
  return (
    <div className="App">
      <Routes
        history={props.history}
        location={props.location}
        match={props.match}
      />
    </div>
  );
}

export default App;
