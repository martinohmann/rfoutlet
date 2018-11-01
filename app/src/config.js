import { cyan500, pink500 } from 'material-ui/styles/colors';

export const strings = {
  appTitle: "rfoutlet",
  on: "On",
  off: "Off",
  toggle: "Toggle",
};

export const styles = {
  appBar: {
    position: "fixed",
    top: 0,
  },
  appBarIcon: {
    display: "none",
  },
  outletGroupContainer: {
    marginTop: 64,
  },
  toolbar: {
    paddingLeft: 16,
    paddingRight: 6,
  },
  outletGroupTitle: {
    fontSize: 14,
  },
  outletGroupButton: {
    margin: 0,
    minWidth: 0,
  },
  outletGroupLabel: {
    padding: 10,
    fontSize: 12,
  },
  buttonOn: {
    color: cyan500,
  },
  buttonOff: {
    color: pink500,
  },
};

export const api = {
  baseUri: `http://${window.location.hostname}:${window.location.port}/api`,
};
