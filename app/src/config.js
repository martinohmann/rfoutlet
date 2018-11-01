import { red500, green500 } from 'material-ui/styles/colors';

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
    color: green500,
  },
  buttonOff: {
    color: red500,
  },
};

export const api = {
  baseUri: `/api`,
};
