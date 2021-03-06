import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import ConfigProvider from '@libs/config/Provider';
import { ThemeProvider } from '@material-ui/styles';
import { createMuiTheme, responsiveFontSizes } from '@material-ui/core/styles';
import App from './features/App';
import * as serviceWorker from './serviceWorker';

const theme = responsiveFontSizes(createMuiTheme());
const jsx = (
  <BrowserRouter>
    <ThemeProvider theme={theme}>
      <ConfigProvider>
        <App />
      </ConfigProvider>
    </ThemeProvider>
  </BrowserRouter>
);

ReactDOM.render(jsx, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
