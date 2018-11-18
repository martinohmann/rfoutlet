import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { TimePicker } from 'material-ui-pickers';

const styles = theme => ({
  timePicker: {
    position: 'fixed',
    top: -1000,
  },
});

class IntervalTimePicker extends React.Component {
  pickerRef = null;

  open = e => this.pickerRef.open(e)

  render() {
    const { classes, value, onChange } = this.props;

    return (
      <TimePicker
        ref={ref => this.pickerRef = ref}
        className={classes.timePicker}
        clearable
        ampm={false}
        value={value}
        onChange={onChange}
      />
    );
  }
}

IntervalTimePicker.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(IntervalTimePicker);
