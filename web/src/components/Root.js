import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { List, NoItemsListItem } from './List';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';

import GithubLink from './GithubLink';
import GroupList from './GroupList';
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
}));

export default function Root() {
  const [groups, setGroups] = useState([]);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    websocket.onMessage(groups => {
      setGroups(groups);
      setLoaded(true);
    });
    websocket.sendMessage({ type: 'status' });
  }, []);

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <AppBar position="fixed">
        <Toolbar>
          <Typography variant="h6" className={classes.title}>
            {config.project.name}
          </Typography>
          <GithubLink url={config.project.url} />
        </Toolbar>
      </AppBar>
      {loaded ? (
        <GroupList groups={groups} />
      ) : (
        <List>
          <NoItemsListItem
            primary="Please wait."
            secondary="Loading outlet states..."
          />
        </List>
      )}
    </div>
  );
}
