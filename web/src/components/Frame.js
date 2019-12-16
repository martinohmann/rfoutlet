import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import SettingsIcon from '@material-ui/icons/Settings';
import config from '../config';
import { useHistory } from 'react-router-dom';

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

export default function Frame({ children }) {
  const classes = useStyles();
  const history = useHistory();

  return (
    <div className={classes.root}>
      <AppBar position="fixed">
        <Toolbar>
          <Typography variant="h6" className={classes.title}>
            {config.project.name}
          </Typography>
          <IconButton className={classes.settings} onClick={() => history.push('/settings')}>
            <SettingsIcon />
          </IconButton>
        </Toolbar>
      </AppBar>
      {children}
    </div>
  );
}
