import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import PowerIcon from '@material-ui/icons/Power';
import PowerOffIcon from '@material-ui/icons/PowerOff';
import SwapHorizIcon from '@material-ui/icons/SwapHoriz';

const styles = theme => ({
  container: {
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
    color: theme.palette.primary[700],
  },
  buttonOff: {
    color: theme.palette.secondary.light,
  },
});

class GroupListItem extends React.Component {
  render() {
    const { classes, name, onActionOn, onActionOff, onActionToggle } = this.props;

    return (
      <ListItem className={classes.container}>
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
    );
  }
}

GroupListItem.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(GroupListItem);
