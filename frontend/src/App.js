import React from 'react';
import darkBaseTheme from 'material-ui/styles/baseThemes/lightBaseTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import AppBar from 'material-ui/AppBar';
import OutletGroupContainer from './OutletGroupContainer';

const config = {
  outlet_groups: [
    {
      identifier: "Group 1",
      outlets: [
        {
          identifier: "Outlet1",
          state: 0,
        },
        {
          identifier: "Outlet2",
          state: 0,
        },
      ], 
    },
    {
      identifier: "Grp 2",
      outlets: [
        {
          identifier: "o1",
          state: 1,
        },
        {
          identifier: "o2",
          state: 0,
        },
      ], 
    },
    {
      identifier: "Group 1",
      outlets: [
        {
          identifier: "Outlet1",
          state: 0,
        },
        {
          identifier: "Outlet2",
          state: 0,
        },
      ], 
    },
    {
      identifier: "Grp 2",
      outlets: [
        {
          identifier: "o1",
          state: 1,
        },
        {
          identifier: "o2",
          state: 0,
        },
      ], 
    },
    {
      identifier: "Group 1",
      outlets: [
        {
          identifier: "Outlet1",
          state: 0,
        },
        {
          identifier: "Outlet2",
          state: 0,
        },
      ], 
    },
    {
      identifier: "Grp 2",
      outlets: [
        {
          identifier: "o1",
          state: 1,
        },
        {
          identifier: "o2",
          state: 0,
        },
      ], 
    },
  ],
}

const styles = {
  appBar: {
    position: "fixed",
    top: 0,
  },
  appBarIcon: {
    display: "none",
  },
  outletGroupContainer: {
    marginTop: 64,
  }
}

const appTitle = "RF-Outlet"

const App = () => (
  <MuiThemeProvider muiTheme={getMuiTheme(darkBaseTheme)}>
    <div>
      <AppBar
        title={appTitle}
        iconStyleLeft={styles.appBarIcon}
        style={styles.appBar}
      />
      <OutletGroupContainer
        outletGroups={config.outlet_groups}
        style={styles.outletGroupContainer}
      />
    </div>
  </MuiThemeProvider>
);

export default App;
