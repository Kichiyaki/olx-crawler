import React from 'react';
import { Switch, Route } from 'react-router-dom';
import MainPage from './MainPage/MainPage';
import ErrorPage from './ErrorPage/ErrorPage';

function App() {
  return (
    <Switch>
      <Route path="/" exact>
        <MainPage />
      </Route>
      <Route path="*">
        <ErrorPage />
      </Route>
    </Switch>
  );
}

export default App;
