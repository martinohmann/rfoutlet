import React from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Switch from '@material-ui/core/Switch';
import EditIcon from '@material-ui/icons/Edit';

import ScheduleDialog from './ScheduleDialog';
import { scheduleToApp } from '../util';

class OutletListItem extends React.Component {
  state = {
    schedule: [],
    scheduleDialogOpen: false,
  }

  static getDerivedStateFromProps(props, state) {
    if (props.schedule !== state.schedule) {
      return {
        schedule: scheduleToApp(props.schedule),
      };
    }

    return null;
  }

  handleToggle = () => {
    const { id } = this.props;

    this.props.dispatchMessage({ type: 'outlet', data: { id, action: 'toggle' } });
  }

  handleScheduleDialogOpen = open => () => {
    this.setState({ scheduleDialogOpen: open });
  }

  render() {
    const { id, name, state } = this.props;
    const { schedule, scheduleDialogOpen } = this.state;

    return (
      <ListItem>
        <ListItemText primary={name} secondary={getScheduleText(schedule)} />
        <ScheduleDialog
          {...this.props}
          outletId={id}
          schedule={schedule}
          open={scheduleDialogOpen}
          onClose={this.handleScheduleDialogOpen(false)}
        />
        <ListItemSecondaryAction>
          <Switch
            color="primary"
            onChange={this.handleToggle}
            checked={state === 1}
            disabled={hasEnabledIntervals(schedule)}
          />
          <IconButton onClick={this.handleScheduleDialogOpen(true)}>
            <EditIcon />
          </IconButton>
        </ListItemSecondaryAction>
      </ListItem>
    );
  }
}

function hasEnabledIntervals(schedule) {
  return schedule.some(interval => interval.enabled)
}

function getScheduleText(schedule) {
  const intervals = schedule.filter(interval => interval.enabled);

  if (intervals.length === 0) {
    return ''
  }

  if (intervals.length === 1) {
    return `1 interval scheduled`
  }

  return `${intervals.length} intervals scheduled`
}

export default OutletListItem;
