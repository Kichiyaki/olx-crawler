import React from 'react';
import { Switch, Route } from 'react-router-dom';
import ConfigPage from './ConfigPage/ConfigPage';
import ErrorPage from './ErrorPage/ErrorPage';
import MainPage from './MainPage/MainPage';
import ObservationsPage from './ObservationsPage/ObservationsPage';

function App() {
  return (
    <Switch>
      <Route path="/" exact>
        <MainPage />
      </Route>
      <Route path="/observations" exact>
        <ObservationsPage />
      </Route>
      <Route path="/config" exact>
        <ConfigPage />
      </Route>
      <Route path="*">
        <ErrorPage />
      </Route>
    </Switch>
  );
}

export default App;
