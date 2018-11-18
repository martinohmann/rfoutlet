import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';

const styles = theme => ({
  title: {
    flexGrow: 1,
    color: theme.palette.common.white,
  },
  toolbarButton: {
    color: theme.palette.common.white,
  },
  container: {
    marginTop: 64,
  },
});

class DialogAppBar extends React.Component {

  render() {
    const {
      classes,
      title,
      onClose,
      onDone,
      doneButtonDisabled,
      doneButtonText
    } = this.props;

    return (
      <AppBar position="fixed">
        <Toolbar>
          <IconButton onClick={onClose} className={classes.toolbarButton}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h6" className={classes.title}>
            {title}
          </Typography>
          {onDone !== undefined ? (
            <Button disabled={doneButtonDisabled} onClick={onDone} className={classes.toolbarButton}>
              {doneButtonText}
            </Button>
          ) : null}
        </Toolbar>
      </AppBar>
    );
  }
}

DialogAppBar.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(DialogAppBar);
