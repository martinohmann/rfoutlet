import React from 'react';
import PropTypes from 'prop-types';

import ConfigurationDialog from './ConfigurationDialog';
import IntervalOptionsList from './IntervalOptionsList';
import IntervalTimePicker from './IntervalTimePicker';
import WeekdaysDialog from './WeekdaysDialog';

class IntervalDialog extends React.Component {
  state = {
    weekdays: [],
    from: null,
    to: null,
    weekdaysDialogOpen: false,
    fromOpen: false,
    toOpen: false,
  }

  static getDerivedStateFromProps(props, state) {
    if (props.id) {
      return { 
        from: props.from,
        to: props.to,
        weekdays: props.weekdays,
      };
    }

    return {
      from: state.from,
      to: state.to,
      weekdays: state.weekdays,
    };
  }

  isIntervalValid() {
    const { from, to, weekdays } = this.state;

    return from !== null && to !== null && weekdays.length > 0;
  }

  handlePickerOpen = (name, open) => e => {
    this.setState({ [name + 'Open']: open });
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
    const { id, enabled, onClose } = this.props;
    const { weekdays, from, to } = this.state;
    const interval = { id, enabled, weekdays, from, to };

    if (id) {
      this.props.onIntervalUpdate(interval);
    } else {
      this.props.onIntervalCreate(interval);
    }

    onClose();
  }

  render() {
    const { open, onClose, id } = this.props;
    const { from, fromOpen, to, toOpen, weekdays, weekdaysDialogOpen } = this.state;

    return (
      <ConfigurationDialog
        title={id ? 'Edit Interval' : 'Add Interval'}
        open={open}
        onClose={onClose}
        onDone={this.handleApply}
        doneButtonDisabled={!this.isIntervalValid()}
        doneButtonText="Apply"
      >
        <IntervalOptionsList
          weekdays={weekdays}
          fromDayTime={from}
          toDayTime={to}
          onWeekdaysClick={this.handleWeekdaysDialogOpen(true)}
          onFromDayTimeClick={this.handlePickerOpen('from', true)}
          onToDayTimeClick={this.handlePickerOpen('to', true)}
        />
        <IntervalTimePicker
          open={fromOpen}
          value={from}
          onChange={this.handlePickerChange('from')}
          onClose={this.handlePickerOpen('from', false)}
        />
        <IntervalTimePicker
          open={toOpen}
          value={to}
          onChange={this.handlePickerChange('to')}
          onClose={this.handlePickerOpen('to', false)}
        />
        <WeekdaysDialog
          open={weekdaysDialogOpen}
          onClose={this.handleWeekdaysDialogOpen(false)}
          onDone={this.handleWeekdaysDialogDone}
          selected={weekdays}
        />
      </ConfigurationDialog>
    );
  }
}

IntervalDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
};

export default IntervalDialog;
