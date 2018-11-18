import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import { TimePicker } from 'material-ui-pickers';

import DialogAppBar from './DialogAppBar';
import WeekdaysDialog from './WeekdaysDialog';
import { formatTime, weekdaysShort } from '../util';

const styles = theme => ({
  container: {
    marginTop: 64,
  },
  timePicker: {
    position: 'fixed',
    top: -1000,
  },
  weekdaySelect: {
    flexGrow: 1,
  },
});

class EditIntervalDialog extends React.Component {
  state = {
    open: false,
    enabled: false,
    create: false,
    weekdays: [],
    from: null,
    to: null,
    weekdaySelectOpen: false,
  }

  pickerRefs = {
    from: null,
    to: null,
  }

  componentWillReceiveProps(nextProps) {
    const { open, id } = nextProps;
    const create = undefined === id;

    console.log(nextProps);
    if (!create) {
      const { from, to, weekdays, enabled } = nextProps;

      this.setState({ from, to, weekdays, enabled });
    }

    this.setState({ open, create });
  }

  intervalValid() {
    const { from, to, weekdays } = this.state;

    return from !== null && to !== null && weekdays.length > 0;
  }

  handleChange = name => date => {
    this.setState({ [name]: date });
  }

  handlePickerOpen = name => e => {
    this.pickerRefs[name].open(e)
  }

  handleWeekdaySelectOpen = open => () => {
    this.setState({ weekdaySelectOpen: open });
  }

  handleWeekdaySelectDone = weekdays => {
    this.setState({ weekdays, weekdaySelectOpen: false });
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

  renderTime(time) {
    if (null === time) {
      return 'unset';
    }

    return formatTime(time);
  }

  render() {
    const { classes, onClose } = this.props;
    const { create, open, weekdays, from, to, weekdaySelectOpen } = this.state;

    return (
      <Dialog fullScreen open={open} onClose={onClose}>
        <DialogAppBar
          title={create ? 'Add Interval' : 'Edit Interval'}
          onClose={onClose}
          onDone={this.handleApply}
          doneButtonDisabled={!this.intervalValid()}
          doneButtonText="Apply"
        />
        <List component="nav" className={classes.container}>
          <ListItem onClick={this.handleWeekdaySelectOpen(true)}>
            <ListItemText
              primary="Weekdays"
              secondary={weekdays.length === 0 ? "none" : weekdays.map(i => weekdaysShort[i]).join(', ')}
            />
          </ListItem>
          <Divider />
          <ListItem onClick={this.handlePickerOpen('from')}>
            <ListItemText primary="From" secondary={this.renderTime(from)} />
          </ListItem>
          <Divider />
          <ListItem onClick={this.handlePickerOpen('to')}>
            <ListItemText primary="To" secondary={this.renderTime(to)} />
          </ListItem>
        </List>
        {this.renderTimePicker('from', from)}
        {this.renderTimePicker('to', to)}
        <WeekdaysDialog
          open={weekdaySelectOpen}
          onClose={this.handleWeekdaySelectOpen(false)}
          onDone={this.handleWeekdaySelectDone}
          selected={weekdays}
        />
      </Dialog>
    );
  }

  renderTimePicker(name, value) {
    const { classes } = this.props;

    return (
      <TimePicker
        ref={ref => this.pickerRefs[name] = ref}
        className={classes.timePicker}
        clearable
        ampm={false}
        value={value}
        onChange={this.handleChange(name)}
      />
    );
  }
}

EditIntervalDialog.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(EditIntervalDialog);
