import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';

import DialogAppBar from './DialogAppBar';
import IntervalOptionsList from './IntervalOptionsList';
import IntervalTimePicker from './IntervalTimePicker';
import WeekdaysDialog from './WeekdaysDialog';

const styles = theme => ({
  container: {
    marginTop: 64,
  },
});

class IntervalDialog extends React.Component {
  state = {
    open: false,
    enabled: false,
    create: false,
    weekdays: [],
    from: null,
    to: null,
    weekdaysDialogOpen: false,
  }

  pickerRefs = {
    from: null,
    to: null,
  }

  componentWillReceiveProps(nextProps) {
    const { open, id } = nextProps;
    const create = undefined === id;

    if (!create) {
      const { from, to, weekdays, enabled } = nextProps;

      this.setState({ from, to, weekdays, enabled });
    }

    this.setState({ open, create });
  }

  isIntervalValid() {
    const { from, to, weekdays } = this.state;

    return from !== null && to !== null && weekdays.length > 0;
  }

  handlePickerOpen = name => e => {
    this.pickerRefs[name].open(e)
  }

  handlePickerChange = name => date => {
    this.setState({ [name]: date });
  }

  handleWeekdaysDialogOpen = open => () => {
    this.setState({ weekdaysDialogOpen: open });
  }

  handleWeekdaysDialogDone = weekdays => {
    this.setState({ weekdays, weekdaysDialogOpen: false });
  }

  handleApply = () => {
    const { id } = this.props;
    const { create, enabled, weekdays, from, to } = this.state;
    const interval = { id, enabled, weekdays, from, to };

    if (create) {
      this.props.onIntervalCreate(interval);
    } else {
      this.props.onIntervalUpdate(interval);
    }

    this.props.onClose();
  }

  render() {
    const { classes, onClose } = this.props;
    const { create, open, from, to, weekdays, weekdaysDialogOpen } = this.state;

    return (
      <Dialog fullScreen open={open} onClose={onClose}>
        <DialogAppBar
          title={create ? 'Add Interval' : 'Edit Interval'}
          onClose={onClose}
          onDone={this.handleApply}
          doneButtonDisabled={!this.isIntervalValid()}
          doneButtonText="Apply"
        />
        <IntervalOptionsList
          className={classes.container}
          weekdays={weekdays}
          fromDayTime={from}
          toDayTime={to}
          onWeekdaysClick={this.handleWeekdaysDialogOpen(true)}
          onFromDayTimeClick={this.handlePickerOpen('from')}
          onToDayTimeClick={this.handlePickerOpen('to')}
        />
        <IntervalTimePicker
          innerRef={ref => this.pickerRefs.from = ref}
          value={from}
          onChange={this.handlePickerChange('from')}
        />
        <IntervalTimePicker
          innerRef={ref => this.pickerRefs.to = ref}
          value={to}
          onChange={this.handlePickerChange('to')}
        />
        <WeekdaysDialog
          open={weekdaysDialogOpen}
          onClose={this.handleWeekdaysDialogOpen(false)}
          onDone={this.handleWeekdaysDialogDone}
          selected={weekdays}
        />
      </Dialog>
    );
  }
}

IntervalDialog.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(IntervalDialog);
