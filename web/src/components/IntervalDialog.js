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
    fromOpen: false,
    toOpen: false,
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
    this.setState({ [name + 'Open']: true });
  }

  handlePickerClose = name => e => {
    this.setState({ [name + 'Open']: false });
  }

  handlePickerChange = name => date => {
    this.setState({ [name]: date, [name + 'Open']: false });
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
    const { create, open, from, fromOpen, to, toOpen, weekdays, weekdaysDialogOpen } = this.state;

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
          open={fromOpen}
          value={from}
          onChange={this.handlePickerChange('from')}
          onClose={this.handlePickerClose('from')}
        />
        <IntervalTimePicker
          open={toOpen}
          value={to}
          onChange={this.handlePickerChange('to')}
          onClose={this.handlePickerClose('to')}
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
