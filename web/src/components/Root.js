import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
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

  useEffect(() => {
    websocket.sendMessage({ type: 'status' });
    websocket.onMessage(setGroups);
  }, [])

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
      <GroupList groups={groups} />
    </div>
  );
}
