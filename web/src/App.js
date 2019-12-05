import React from 'react';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import LuxonUtils from '@date-io/luxon';
import cyan from '@material-ui/core/colors/cyan';
import { HashRouter as Router } from 'react-router-dom';
import { I18nextProvider } from 'react-i18next';

import Root from './components/Root';
import i18n from './i18n';

const theme = createMuiTheme({
  palette: {
    primary: cyan,
  },
  typography: {
    useNextVariants: true,
  },
  overrides: {
    MuiPickersToolbarText: {
      toolbarTxt: {
        color: "rgba(255, 255, 255, 0.54)",
      },
      toolbarBtnSelected: {
        color: "white",
      },
    },
  },
});

export default function App() {
  return (
    <Router>
      <I18nextProvider i18n={i18n}>
        <MuiThemeProvider theme={theme}>
          <MuiPickersUtilsProvider utils={LuxonUtils} locale={i18n.language}>
            <Root />
          </MuiPickersUtilsProvider>
        </MuiThemeProvider>
      </I18nextProvider>
    </Router>
  );
}
