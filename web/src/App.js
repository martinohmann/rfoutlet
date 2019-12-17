import React, { useState, useEffect } from 'react';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import LuxonUtils from '@date-io/luxon';
import cyan from '@material-ui/core/colors/cyan';
import { HashRouter } from 'react-router-dom';
import { I18nextProvider } from 'react-i18next';
import { GroupProvider } from './Context';
import i18n from './i18n';
import Routes from './components/Routes';
import dispatcher from './dispatcher';

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
  const [groups, setGroups] = useState([]);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    dispatcher.addMessageListener(groups => {
      setGroups(groups);
      setReady(true);
    });

    dispatcher.dispatchStatusMessage();
  }, []);

  return (
    <HashRouter>
      <I18nextProvider i18n={i18n}>
        <MuiThemeProvider theme={theme}>
          <MuiPickersUtilsProvider utils={LuxonUtils} locale={i18n.language}>
            <GroupProvider groups={groups}>
              <Routes ready={ready} groups={groups} />
            </GroupProvider>
          </MuiPickersUtilsProvider>
        </MuiThemeProvider>
      </I18nextProvider>
    </HashRouter>
  );
}
