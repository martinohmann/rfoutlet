import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import { TimePicker } from '@material-ui/pickers';

const useStyles = makeStyles({
  timePicker: {
    position: 'fixed',
    top: -1000,
  },
});

export default function IntervalTimePicker(props) {
  const classes = useStyles();

  return (
    <TimePicker
      className={classes.timePicker}
      clearable
      ampm={false}
      {...props}
    />
  );
}

IntervalTimePicker.propTypes = {
  open: PropTypes.bool.isRequired,
  onChange: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};
