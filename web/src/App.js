import React from 'react';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import LuxonUtils from '@date-io/luxon';
import cyan from '@material-ui/core/colors/cyan';

import Root from './components/Root';

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
    <MuiThemeProvider theme={theme}>
      <MuiPickersUtilsProvider utils={LuxonUtils}>
        <Root />
      </MuiPickersUtilsProvider>
    </MuiThemeProvider>
  );
}
