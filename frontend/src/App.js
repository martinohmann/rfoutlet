import React from 'react';
import lightBaseTheme from 'material-ui/styles/baseThemes/lightBaseTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import AppBar from 'material-ui/AppBar';
import OutletGroupContainer from './components/OutletGroupContainer';
import { strings, styles } from './config';

const App = () => (
  <MuiThemeProvider muiTheme={getMuiTheme(lightBaseTheme)}>
    <div>
      <AppBar
        title={strings.appTitle}
        iconStyleLeft={styles.appBarIcon}
        style={styles.appBar}
      />
      <OutletGroupContainer />
    </div>
  </MuiThemeProvider>
);

export default App;
