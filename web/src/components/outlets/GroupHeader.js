import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import PowerIcon from '@material-ui/icons/Power';
import PowerOffIcon from '@material-ui/icons/PowerOff';
import SwapHorizIcon from '@material-ui/icons/SwapHoriz';

const useStyles = makeStyles(theme => ({
  container: {
    padding: 0,
  },
  header: {
    paddingTop: 1,
    paddingBottom: 1,
    paddingRight: 6,
    background: theme.palette.grey[100],
  },
  groupName: {
    flexGrow: 1,
    fontWeight: 700,
    color: theme.palette.grey[800],
  },
  buttonOn: {
    color: theme.palette.primary.dark,
  },
  buttonOff: {
    color: theme.palette.secondary.light,
  },
}));

export default function GroupHeader(props) {
  const classes = useStyles();

  const { name, onActionOn, onActionOff, onActionToggle } = props;

  return (
    <List className={classes.container}>
      <ListItem className={classes.header}>
        <ListItemText className={classes.groupName} primary={name} disableTypography={true} />
        <IconButton className={classes.buttonOff} onClick={onActionOff}>
          <PowerOffIcon />
        </IconButton>
        <IconButton className={classes.buttonOn} onClick={onActionOn}>
          <PowerIcon />
        </IconButton>
        <IconButton onClick={onActionToggle}>
          <SwapHorizIcon />
        </IconButton>
      </ListItem>
    </List>
  );
}

GroupHeader.propTypes = {
  name: PropTypes.string.isRequired,
  onActionOn: PropTypes.func.isRequired,
  onActionOff: PropTypes.func.isRequired,
  onActionToggle: PropTypes.func.isRequired,
};
