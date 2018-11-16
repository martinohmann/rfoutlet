import React from 'react';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import { MuiPickersUtilsProvider } from 'material-ui-pickers';
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
});

export default class App extends React.Component {
  render() {
    return (
      <MuiThemeProvider theme={theme}>
        <MuiPickersUtilsProvider utils={LuxonUtils}>
          <Root />
        </MuiPickersUtilsProvider>
      </MuiThemeProvider>
    );
  }
}
