import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { List, NoItemsListItem } from './List';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { useTranslation } from 'react-i18next';
import IconButton from '@material-ui/core/IconButton';
import SettingsIcon from '@material-ui/icons/Settings';
import { Route, Switch, useHistory } from 'react-router-dom';

import GroupList from './GroupList';
import SettingsDialog from './SettingsDialog';
import config from '../config';
import websocket from '../websocket';

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },
  title: {
    flexGrow: 1,
    color: theme.palette.common.white,
  },
  settings: {
    color: theme.palette.common.white,
  },
}));

export default function Root() {
  const [groups, setGroups] = useState([]);
  const [loaded, setLoaded] = useState(false);
  const history = useHistory();

  useEffect(() => {
    websocket.onMessage(groups => {
      setGroups(groups);
      setLoaded(true);
    });
    websocket.sendMessage({ type: 'status' });
  }, []);


  const handleDialogOpen = (open) => () => history.push(open ? '/settings' : '/');

  const classes = useStyles();
  const { t } = useTranslation();

  return (
    <div className={classes.root}>
      <AppBar position="fixed">
        <Toolbar>
          <Typography variant="h6" className={classes.title}>
            {config.project.name}
          </Typography>
          <IconButton className={classes.settings} onClick={handleDialogOpen(true)}>
            <SettingsIcon />
          </IconButton>
        </Toolbar>
      </AppBar>
      <Switch>
        <Route path="/settings">
          <SettingsDialog onClose={handleDialogOpen(false)} />
        </Route>
        <Route path="/">
          {loaded ? (
            <GroupList groups={groups} />
          ) : (
            <List>
              <NoItemsListItem
                primary={t('loading-primary')}
                secondary={t('loading-secondary')}
              />
            </List>
          )}
        </Route>
      </Switch>
    </div>
  );
}
