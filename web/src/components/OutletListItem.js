import React from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Switch from '@material-ui/core/Switch';
import EditIcon from '@material-ui/icons/Edit';

import ScheduleDialog from './ScheduleDialog';
import { scheduleToApp, formatSchedule } from '../util';

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
    const { id, dispatchMessage } = this.props;

    dispatchMessage({ type: 'outlet', data: { id, action: 'toggle' } });
  }

  handleScheduleDialogOpen = open => () => {
    console.log(open);
    this.setState({ scheduleDialogOpen: open });
  }

  render() {
    const { id, name, state } = this.props;
    const { schedule, scheduleDialogOpen } = this.state;

    return (
      <ListItem>
        <ListItemText primary={name} secondary={formatSchedule(schedule)} />
        <ScheduleDialog
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
            disabled={schedule.some(interval => interval.enabled)}
          />
          <IconButton onClick={this.handleScheduleDialogOpen(true)}>
            <EditIcon />
          </IconButton>
        </ListItemSecondaryAction>
      </ListItem>
    );
  }
}

export default OutletListItem;
