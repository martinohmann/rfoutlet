import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';

const useStyles = makeStyles(theme => ({
  title: {
    flexGrow: 1,
    color: theme.palette.common.white,
  },
  toolbar: {
    paddingLeft: theme.spacing(1),
    paddingRight: theme.spacing(1),
  },
  toolbarButton: {
    color: theme.palette.common.white,
  },
  container: {
    marginTop: 64,
  },
}));

export default function DialogAppBar(props) {
  const classes = useStyles();

  const {
    title,
    onClose,
    onDone,
    doneButtonDisabled,
    doneButtonText
  } = props;

  return (
    <AppBar position="fixed">
      <Toolbar className={classes.toolbar}>
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

DialogAppBar.propTypes = {
  onClose: PropTypes.func.isRequired,
  onDone: PropTypes.func,
  title: PropTypes.string.isRequired,
};
